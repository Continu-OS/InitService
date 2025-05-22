package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/Continu-OS/InitService/src/handler"
	"github.com/Continu-OS/InitService/src/linux"
	"github.com/Continu-OS/InitService/src/servcaprocman"
)

func main() {
	// Es wird geprüft ob das Programm als Init sowie mit Rootrechten ausgeführt wird
	runAsInitRoot, err := linux.RunAsInitProcess()
	if err != nil {
		fmt.Println("Error checking if running as init process:", err)
		os.Exit(1)
	}
	if !runAsInitRoot {
		fmt.Println("Not running as init service with root privileges")
		os.Exit(1)
	}

	// Der SIGNAL Handler wird gestartet
	if err := handler.StartHostKernelSIGNALHandler(); err != nil {
		fmt.Println("Failed to start host kernel SIGNAL handler:", err)
		os.Exit(1)
	}

	// Der Zombie Prozess Handler wird gestartet
	if handler.StartReapZombieProcessesHandler(); err != nil {
		fmt.Println("Failed to start host zombie process handler:", err)
		os.Exit(1)
	}

	// Die Bootloader Parameter werden ausgelesen
	discoveredBootParameters, err := linux.GetAllBootloaderParameters()
	if err != nil {
		fmt.Println("Failed to get bootloader parameters:", err)
		os.Exit(1)
	}

	// Überprüft, ob bei Embedded-Systemen (z. B. Raspberry Pi) die Root-Partition beim Booten angepasst werden muss,
	// um die gesamte Kapazität der SD-Karte, etc... zu nutzen.
	// Dies verhindert Systemfehler durch unvollständige Speicherinitialisierung.
	// Die Funktion ist nur für SoC-Systeme aktiv. Auf Workstations, Laptops, Servern, VMs usw. bleibt sie deaktiviert.

	hardCheckRootSystemDrive := false

	adjustVolume := discoveredBootParameters["check_andor_adjust_rvolume.contios.github.com"]
	embeddedBoard := strings.EqualFold(discoveredBootParameters["environment.contios.github.com"], "EMBEDDED_BOARD")

	if embeddedBoard {
		// Standardmäßig aktivieren, außer es ist explizit auf FALSE gesetzt
		if !strings.EqualFold(adjustVolume, "FALSE") {
			hardCheckRootSystemDrive = true
		}
	} else if strings.EqualFold(adjustVolume, "TRUE") {
		// Bei Nicht-SoC-Systemen nur aktivieren, wenn explizit TRUE gesetzt wurde
		hardCheckRootSystemDrive = true
	}

	afterHardDriveCheckProcessRebootHostSystem := false

	if hardCheckRootSystemDrive {
		// Ermittelt das Root-Gerät
		rootDev, err := linux.GetRootDevice()
		if err != nil {
			panic("Fehler beim Ermitteln des Root-Geräts: " + err.Error())
		}
		// Ermittelt das zugehörige Blockgerät des Root-Partitionsgeräts
		device := linux.GetDeviceFromPartition(rootDev)
		// Prüft, ob eine Größenanpassung der Root-Partition notwendig ist
		needed, err := linux.NeedsRootDriveResize(rootDev, device)
		if err != nil {
			panic("Fehler bei der Prüfung, ob eine Größenanpassung nötig ist: " + err.Error())
		}

		if needed {
			// Führt die Größenanpassung durch, falls notwendig
			if err := linux.ResizeRootDriveIfNeeded(device, rootDev); err != nil {
				panic("Fehler beim Anpassen der Root-Partition: " + err.Error())
			}

			// Merkt, dass nach dem Vorgang ein Neustart des Systems erforderlich ist
			afterHardDriveCheckProcessRebootHostSystem = true
		}
	}

	// Es wird geprüft ob der Computer neugestartet werden muss
	if afterHardDriveCheckProcessRebootHostSystem {
		if err := linux.RebootSystemNow(); err != nil {
			panic(err)
		}
	}

	// Die Einstellungen werden geladen

	// Der Dienst und Prozess Manager wird gestartet
	if err := servcaprocman.Init(); err != nil {
		fmt.Println("Failed to Init Process and Services Manager:", err)
		os.Exit(1)
	}

	// Die Sekundären Init Service Dienste werden gestartet, diese Dienste sind
	// festerbestandteil des Betrybsystem und können nicht durch Einstellungen in /System/Library,
	// oder in /Library/ verändert werden. Dazu zählen unteranderen:
	// -- IPC Service
	// -- Wechseldatenträgerverwaltung
	// -- Netzwerkverwaltung
	// Diese Dienste sind fester bestandteil des ContOS InitServices

}

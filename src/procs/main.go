package main

import (
	"fmt"
	"log"
	"os"
	"strings"

	InitService "github.com/Continu-OS/InitService/src"
	"github.com/Continu-OS/InitService/src/cformat"
	"github.com/Continu-OS/InitService/src/config"
	"github.com/Continu-OS/InitService/src/handler"
	"github.com/Continu-OS/InitService/src/linux"
	"github.com/Continu-OS/InitService/src/linux/firmware"
	"github.com/Continu-OS/InitService/src/servcaprocman"
)

func fatalClose(msg ...any) {
	fmt.Println(msg...)
	cformat.ConsoleExit()
	os.Exit(1)
}

func main() {
	// Der Konsolenlog wird Preparier
	cformat.ConsoleEntering()
	defer cformat.ConsoleExit()

	// Es wird geprüft ob das Programm als Init sowie mit Rootrechten ausgeführt wird
	/*runAsInitRoot, err := linux.RunAsInitProcess()
	if err != nil {
		fmt.Println("Error checking if running as init process:", err)
		os.Exit(1)
	}
	if !runAsInitRoot {
		fmt.Println("Not running as init service with root privileges")
		os.Exit(1)
	}
	*/

	// Log
	log.Println("ContinuOS Init Service started")

	// Der SIGNAL Handler wird gestartet
	if err := handler.StartHostKernelSIGNALHandler(); err != nil {
		fatalClose("Failed to start host kernel SIGNAL handler:", err.Error())
	}

	// Der Zombie Prozess Handler wird gestartet
	if err := handler.StartReapZombieProcessesHandler(); err != nil {
		fatalClose("Failed to start host zombie process handler:", err.Error())
	}

	// Die Bootloader Parameter werden ausgelesen
	bootArgs, err := firmware.GetAllBootloaderParameters()
	if err != nil {
		fatalClose("Failed to get bootloader parameters:", err.Error())
	}

	// Ermittelt das Root-Gerät
	rootPart, err := linux.GetRootDevicePartition()
	if err != nil {
		fatalClose("Fehler beim Ermitteln des Root-Geräts: ", err.Error())
	}

	// Log
	log.Printf("System Partition: %s\n", rootPart)

	// Ermittelt das zugehörige Blockgerät des Root-Partitionsgeräts
	rootDevice := linux.GetDeviceFromPartition(InitService.MemoryPartition(rootPart))

	// Log
	log.Printf("System Rootdevice: %s\n", rootDevice)

	// Überprüft, ob bei Embedded-Systemen (z. B. Raspberry Pi) die Root-Partition beim Booten angepasst werden muss,
	// um die gesamte Kapazität der SD-Karte, etc... zu nutzen.
	// Dies verhindert Systemfehler durch unvollständige Speicherinitialisierung.
	// Die Funktion ist nur für SoC-Systeme aktiv. Auf Workstations, Laptops, Servern, VMs usw. bleibt sie deaktiviert.

	hardCheckRootSystemDrive := false

	adjustVolume := bootArgs["check_andor_adjust_rvolume.contios.github.com"]
	hostIsSBCS := strings.EqualFold(bootArgs["environment.contios.github.com"], "SBCS")

	switch {
	case hostIsSBCS && adjustVolume != "FALSE":
		hardCheckRootSystemDrive = true
	case !hostIsSBCS && adjustVolume == "TRUE":
		hardCheckRootSystemDrive = true
	}

	afterHardDriveCheckProcessRebootHostSystem := false

	if hardCheckRootSystemDrive {
		// Prüft, ob eine Größenanpassung der Root-Partition notwendig ist
		needed, err := linux.NeedsRootDriveResize(InitService.RootDevice(rootDevice), InitService.RootPartition(rootPart))
		if err != nil {
			fatalClose("Fehler bei der Prüfung, ob eine Größenanpassung nötig ist: ", err.Error())
		}

		if needed {
			// Führt die Größenanpassung durch, falls notwendig
			if err := linux.ResizeRootDriveIfNeeded(rootDevice, InitService.MemoryPartition(rootPart)); err != nil {
				fatalClose("Fehler beim Anpassen der Root-Partition: ", err.Error())
			}

			// Merkt, dass nach dem Vorgang ein Neustart des Systems erforderlich ist
			afterHardDriveCheckProcessRebootHostSystem = true
		}
	}

	// Es wird geprüft ob der Computer neugestartet werden muss
	if afterHardDriveCheckProcessRebootHostSystem {
		if err := linux.RebootSystemNow(); err != nil {
			fatalClose(err.Error())
		}
	}

	// !!!!
	// Ab hier ist das System soweit sicher geladen dass keine Fehler am Dateisystem mehr auftreten dürften
	// !!!!

	// Es wird geprüft ob die Festpallte Verifizierbar ist,
	// betrifft nicht die Daten direkt auf der Festplatte, sondern nur das Paritionierungslayout.
	verifyed, err := linux.VerifyRootSystemDriveSignature(rootPart, rootDevice)
	if err != nil {
		fatalClose(err.Error())
	}
	if !verifyed {
		fatalClose("Invlaid Rootdrive Filesystem, rebooting now...")
	}

	// Die Einstellungen werden geladen
	if err := config.LoadHostInitConfig(bootArgs); err != nil {
		fatalClose("Error by loading core settings", err.Error())
	}

	// Der Dienst und Prozess Manager wird gestartet
	if err := servcaprocman.Init(); err != nil {
		fatalClose("Failed to Init Process and Services Manager:", err.Error())
	}

	// Die Sekundären Init Service Dienste werden gestartet, diese Dienste sind
	// festerbestandteil des Betrybsystem und können nicht durch Einstellungen in /System/Library,
	// oder in /Library/ verändert werden. Dazu zählen unteranderen:
	// -- IPC Service
	// -- Wechseldatenträgerverwaltung
	// -- Netzwerkverwaltung
	// Diese Dienste sind fester bestandteil des ContOS InitServices, ohne diese Dienste bricht der Init Vorgang ab
}

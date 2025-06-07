package main

import (
	"fmt"
	"log"
	"os"

	InitService "github.com/Continu-OS/InitService/src/init"
	"github.com/Continu-OS/InitService/src/init/cformat"
	"github.com/Continu-OS/InitService/src/init/config"
	"github.com/Continu-OS/InitService/src/init/defconfig"
	"github.com/Continu-OS/InitService/src/init/frameworks"
	"github.com/Continu-OS/InitService/src/init/fs"
	"github.com/Continu-OS/InitService/src/init/handler"
	"github.com/Continu-OS/InitService/src/init/linux"
	"github.com/Continu-OS/InitService/src/init/linux/cgroups"
	"github.com/Continu-OS/InitService/src/init/linux/firmware"
	"github.com/Continu-OS/InitService/src/init/runtimes"
	"github.com/Continu-OS/InitService/src/init/servcaprocman"
	"github.com/Continu-OS/InitService/src/init/toolsets"
)

func fatalClose(msg ...any) {
	log.Printf("FATAL ERROR: %v", fmt.Sprint(msg...))

	// Versuche sauberen Shutdown wenn möglich
	if err := linux.TryGracefulShutdown(); err != nil {
		log.Printf("Graceful shutdown failed: %v", err)
	}

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

	// Die Defconfig Parameter werden vorbereitet
	if defconfig.InitWithBootArgs(&bootArgs); err != nil {
		fatalClose("Fatal error, can't inital defconfig.... ", err.Error())
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
	// Hierfür muss ermittelt werden in was für einer Umgebung der Dienst ausgeführt wird.
	if defconfig.GetBoolOption("memory.auto_resize_enable") && !cgroups.IsRunningInDocker() {
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

			// Der Host muss neugestartet werden
			if err := linux.RebootSystemNow(); err != nil {
				fatalClose(err.Error())
			}
		}
	} else {
		log.Println("Root auto-resize skipped – either disabled in config or running inside Docker")

	}

	// !!!!
	// Ab hier ist das System soweit sicher geladen dass keine Fehler am Dateisystem mehr auftreten dürften
	// !!!!

	// Die Einstellungen werden geladen
	if err := config.LoadHostInitConfig(); err != nil {
		fatalClose("Error by loading core settings", err.Error())
	}

	// Die IPC-API Runtime wird gestartet
	if err := runtimes.StartIPCRuntime(); err != nil {
		fatalClose("Error by starting IPC-API runtime", err.Error())
	}

	// Die Netzwerk Runtime wird gestartet
	if err := runtimes.StartNetworkRuntime(); err != nil {
		fatalClose("Error by starting Network runtime", err.Error())
	}

	// Die Geräteerkennung wird gestartet, wichtig ist das die Geräteerkennung mit zum "Schluss" gestartet wird.
	// Die Geräteereknnung Signalisiert den Runtimes, sowie den Service Manager bestimmte Events.
	if err := handler.StartDeviceEventHandler(); err != nil {
		fatalClose("Error by starting device handler", err.Error())
	}

	// Die Sekundären Init Service Dienste werden gestartet, diese Dienste sind
	// festerbestandteil des Betrybsystem und können nicht durch Einstellungen in /System/Library,
	// oder in /Library/ verändert werden. Dazu zählen unteranderen:
	// -- Wechseldatenträgerverwaltung
	// -- Netzwerkverwaltung
	// Diese Dienste sind fester bestandteil des ContOS InitServices, ohne diese Dienste bricht der Init Vorgang ab

	// Der ContinOS Core Dienst wird gestartet, dieser Dienst ist Essentziell und Stelle reine reihen an Basisfunktionen bereit
	// diese Basisfunktionen von ContinOS können über die API verwendet werden.
	if err := servcaprocman.RunServiceProcessWithoutContainer(0, 0, fmt.Sprintf("/sbin/cntincore -initSecID=0")); err != nil {
		fatalClose("Error by starting device service", err.Error())
	}

	// Diese Dienste werden nur benötigt wenn das Programm außerhalb von DOCKER ausgeführt wird
	if !cgroups.IsRunningInDocker() {
		// Der Datenträger Dienst wird gestartet
		if err := servcaprocman.RunServiceProcessWithoutContainer(0, 0, fmt.Sprintf("/sbin/medevcn -initSecID=0")); err != nil {
			fatalClose("Error by starting device service", err.Error())
		}

		// Der Netzwerkdienst wird gestartet
		if err := servcaprocman.RunServiceProcessWithoutContainer(0, 0, fmt.Sprintf("/sbin/netwcns -initSecID=0")); err != nil {
			fatalClose("Error by starting network service", err.Error())
		}
	} else {
		log.Println("Docker detected – skipping startup of media and network services")

	}

	// Die System Frameworks werden geladen
	systemBaseFrameworks, err := fs.GetAllBaseSystemFrameworks()
	if err != nil {
		fatalClose("Error by loading Basesystem Frameworks", err.Error())
	}
	for _, baseSystemFrameworkFile := range systemBaseFrameworks {
		// Das Basis Framework wird geladen
		if err := frameworks.InloadBaseFramework(baseSystemFrameworkFile); err != nil {
			fatalClose("Error by loading Base System Framework", err.Error())
		}
	}

	// Es werden alle Verfügbaren Basis System Toolsets geladen
	systemBaseToolsets, err := fs.GetAllBaseSystemToolsets()
	if err != nil {
		fatalClose("Error by starting network service", err.Error())
	}
	for _, baseSystemToolset := range systemBaseToolsets {
		if err := toolsets.InloadBaseToolset(baseSystemToolset); err != nil {
			fatalClose("Error by loading Base System Toolset", err.Error())
		}
	}

	// Es werden alle Verfügbaren Basis System Dienste geladen
	systemBaseServices, err := fs.GetAllBaseSystemServices()
	if err != nil {
		fatalClose("Error by starting network service", err.Error())
	}
	for _, baseSystemService := range systemBaseServices {
		if err := servcaprocman.InloadBaseFramework(baseSystemService); err != nil {
			fatalClose("Error by loading Base System Service", err.Error())
		}
	}
}

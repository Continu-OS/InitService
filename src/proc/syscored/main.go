package main

import (
	"fmt"
	"log"
	"os"

	"github.com/Continu-OS/syscored/src/proc"
	InitService "github.com/Continu-OS/syscored/src/syscored"
	"github.com/Continu-OS/syscored/src/syscored/cformat"
	"github.com/Continu-OS/syscored/src/syscored/config"
	"github.com/Continu-OS/syscored/src/syscored/defconfig"
	"github.com/Continu-OS/syscored/src/syscored/fs"
	"github.com/Continu-OS/syscored/src/syscored/handler"
	"github.com/Continu-OS/syscored/src/syscored/linux"
	"github.com/Continu-OS/syscored/src/syscored/linux/cgroups"
	"github.com/Continu-OS/syscored/src/syscored/linux/firmware"
	"github.com/Continu-OS/syscored/src/syscored/servcaprocman"
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

	// Der Notify (XPC/IPC) Dienst wird gestartet
	if err := servcaprocman.RunServiceSupervised(0, 0, fmt.Sprintf("/sbin/notifyd -initSecID=%s", "0")); err != nil {
		fatalClose("Error by starting notifyd service", err.Error())
	}

	// Der Sicherheitsdienst securityd wird gestartet
	if err := servcaprocman.RunServiceSupervised(0, 0, fmt.Sprintf("/sbin/securityd -initSecID=%s", "0")); err != nil {
		fatalClose("Error by starting securityd service", err.Error())
	}

	// Diese Dienste werden nur benötigt wenn das Programm außerhalb von DOCKER ausgeführt wird
	if !cgroups.IsRunningInDocker() {
		// Der Power Dienst wird gestartet
		if err := servcaprocman.RunServiceSupervised(0, 0, fmt.Sprintf("/sbin/powerd -initSecID=%s", "0")); err != nil {
			fatalClose("Error by starting deviced service", err.Error())
		}

		// Der Geräte Dienst wird gestartet
		if err := servcaprocman.RunServiceSupervised(0, 0, fmt.Sprintf("/sbin/deviced -initSecID=%s", "0")); err != nil {
			fatalClose("Error by starting deviced service", err.Error())
		}

		// Der Netzwerkdienst wird gestartet
		if err := servcaprocman.RunServiceSupervised(0, 0, fmt.Sprintf("/sbin/networkd -initSecID=%s", "0")); err != nil {
			fatalClose("Error by starting networkd service", err.Error())
		}
	} else {
		log.Println("Docker detected – skipping startup of devices and network services")
	}

	// Es werden alle Verfügbaren Basis System Dienste geladen
	systemBaseServices, err := fs.GetAllBaseSystemServices()
	if err != nil {
		fatalClose(err)
	}
	for _, baseSystemService := range systemBaseServices {
		if err := servcaprocman.InloadService(baseSystemService); err != nil {
			fatalClose("Error by loading Base System Service", err.Error())
		}
	}

	// Es wird geprüft ob die Abhänigkeiten für die Dienste aufgelöst werden könnnen
	// Frameworks, Toolsets, andere Dienste
	if servcaprocman.VerifyServicesFrameworkDependencies(); err != nil {
		fatalClose(err)
	}
	if servcaprocman.VerifyServicesToolsetDependencies(); err != nil {
		fatalClose(err)
	}
	if servcaprocman.VerifyServicesServiceDependencies(); err != nil {
		fatalClose(err)
	}

	// Es werden alle Dienst gestartet
	if servcaprocman.StartAllServicesOptimal(); err != nil {
		fatalClose(err)
	}

	// Es wird geprüft ob in der DEF_CONFIG ein Standard Prozess eingetragen ist welcher gestartet werden soll
	if val, hasOption := defconfig.GetOption("system.static_login_promt"); hasOption {
		// Es wird geprüft ob die Datei vorhanden ist
		strPathToLoginPromtTool, IsOk := val.(string)
		if !IsOk {
			fatalClose("Invalid DEFCONFIG")
		}
		if !fs.FileExists(strPathToLoginPromtTool) {
			fatalClose(fmt.Sprintf("Invalid DEFCONFIG, no LOGIN PROMT SERVICE FOUND ON %s", strPathToLoginPromtTool))
		}
	}

	// Der Prozess wird offen gehalten
	proc.RunnerProcs()
}

package handler

import (
	"log"
	"os"
	"os/signal"
	"syscall"
)

func StartHostKernelSIGNALHandler() error {
	signals := make(chan os.Signal, 1)
	signal.Notify(signals,
		syscall.SIGTERM,
		syscall.SIGINT,
		syscall.SIGQUIT,
		syscall.SIGHUP,
		syscall.SIGUSR1,
		syscall.SIGUSR2,
	)

	go func() {
		defer func() {
			log.Println("Kernel Signal handler was stoped")
		}()

		log.Println("Kernel Signal handler was started")

		for sig := range signals {
			log.Printf("Initdienst empfängt Signal: %v", sig)
			switch sig {
			case syscall.SIGTERM, syscall.SIGINT, syscall.SIGQUIT:
				log.Println("System wird sauber heruntergefahren ...")
				// Hier alle Dienste stoppen, fs syncen, Reboot oder Shutdown einleiten
				ShutdownAllServices()
				SyncFilesystems()
				os.Exit(0)
			case syscall.SIGHUP:
				log.Println("Config-Reload (SIGHUP) angefordert")
				ReloadConfiguration()
			case syscall.SIGUSR1:
				log.Println("SIGUSR1: Soft-Restart oder spezielle Aufgabe")
				SoftRestart()
			case syscall.SIGUSR2:
				log.Println("SIGUSR2: Re-Konfiguration oder andere Aufgabe")
				ReloadServices()
			}
		}
	}()
	return nil
}

// Dummyfunktionen für Beispiel:
func ShutdownAllServices() {}
func SyncFilesystems()     {}
func ReloadConfiguration() {}
func SoftRestart()         {}
func ReloadServices()      {}

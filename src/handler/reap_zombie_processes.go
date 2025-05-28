package handler

import (
	"log"
	"os"
	"os/signal"
	"syscall"
)

func StartReapZombieProcessesHandler() error {
	sigchld := make(chan os.Signal, 1)
	signal.Notify(sigchld, syscall.SIGCHLD)

	go func() {
		defer func() {
			log.Println("Zombie Process reaper was started")
		}()

		log.Println("Zombie Process reaper was started")

		for range sigchld {
			for {
				var status syscall.WaitStatus
				var rusage syscall.Rusage
				pid, err := syscall.Wait4(-1, &status, syscall.WNOHANG, &rusage)
				if err != nil {
					if err == syscall.ECHILD {
						break
					}
					log.Printf("Fehler beim Reaping: %v", err)
					break
				}
				if pid <= 0 {
					break
				}
				log.Printf("Reaped process %d: status=%v", pid, status)
			}
		}
	}()
	return nil
}

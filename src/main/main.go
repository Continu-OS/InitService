package main

import (
	"fmt"
	"os"
	"os/user"

	"github.com/Continu-OS/InitService/src/handler"
)

func main() {
	// Prüfen ob PID 1
	pid := os.Getpid()
	if pid == 1 {
		fmt.Println("Läuft als PID 1")
		os.Exit(1)
	} else {
		fmt.Println("Nicht PID 1, sondern:", pid)
		os.Exit(1)
	}

	// Prüfen ob als Root
	u, err := user.Current()
	if err != nil {
		fmt.Println("Fehler beim Abrufen des Benutzers:", err)
		os.Exit(1)
	}

	if u.Uid == "0" {
		fmt.Println("Läuft als Root")
		os.Exit(1)
	} else {
		fmt.Println("Nicht als Root (UID =", u.Uid+")")
		os.Exit(1)
	}

	// Alternative: UID direkt über syscall
	uid := os.Geteuid()
	if uid == 0 {
		fmt.Println("Echt als Root (EUID = 0)")
		os.Exit(1)
	} else {
		fmt.Println("EUID ist", uid)
		os.Exit(1)
	}

	// Der SIGNAL Handler wird gestartet
	if err := handler.StartHostKernelSIGNALHandler(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// Es wird geprüft ob ein Parameter übergeben wurde
}

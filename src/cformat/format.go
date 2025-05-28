package cformat

import (
	"os"
)

// BoldPrintf druckt alles fett formatiert
func ConsoleEntering() {
	os.Stdout.Write([]byte("\033[1m"))
	os.Stderr.Write([]byte("\033[1m"))
}

func ConsoleExit() {
	os.Stdout.Write([]byte("\033[0m"))
	os.Stderr.Write([]byte("\033[0m"))
}

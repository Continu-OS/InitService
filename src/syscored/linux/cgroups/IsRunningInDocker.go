package cgroups

import (
	"bufio"
	"os"
	"strings"
)

// IsRunningInDocker prüft, ob der Prozess in einem Docker-Container läuft.
func IsRunningInDocker() bool {
	// 1. Prüfe auf Docker-spezifische cgroup-Einträge
	if f, err := os.Open("/proc/1/cgroup"); err == nil {
		defer f.Close()
		scanner := bufio.NewScanner(f)
		for scanner.Scan() {
			line := scanner.Text()
			if strings.Contains(line, "docker") || strings.Contains(line, "containerd") {
				return true
			}
		}
	}

	// 2. Prüfe auf die Existenz der Docker-Umgebung
	if _, err := os.Stat("/.dockerenv"); err == nil {
		return true
	}

	// 3. Optional: Weitere heuristische Checks (wenn nötig)

	return false
}

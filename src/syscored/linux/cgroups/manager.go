package cgroups

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
)

const basePath = "/sys/fs/cgroup"

// Create erstellt eine neue Cgroup mit Speicher- und CPU-Limit
func Create(name string, memMB int, cpuPercent int) error {
	path := filepath.Join(basePath, "contios", name)
	if err := os.MkdirAll(path, 0755); err != nil {
		return err
	}
	if err := os.WriteFile(filepath.Join(path, "memory.max"),
		[]byte(fmt.Sprintf("%d", memMB*1024*1024)), 0644); err != nil {
		return err
	}
	if err := os.WriteFile(filepath.Join(path, "cpu.max"),
		[]byte(fmt.Sprintf("%d 100000", cpuPercent*1000)), 0644); err != nil {
		return err
	}
	return nil
}

// Assign fügt einen laufenden Prozess der Cgroup hinzu
func Assign(name string, pid int) error {
	return os.WriteFile(filepath.Join(basePath, "contios", name, "cgroup.procs"),
		[]byte(strconv.Itoa(pid)), 0644)
}

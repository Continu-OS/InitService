package linux

import (
	"regexp"
	"strings"

	InitService "github.com/Continu-OS/syscored/src/syscored"
)

// Extrahiert das Basisspeichergerät aus der Partition, z.B.
// /dev/mmcblk0p1 -> /dev/mmcblk0
// /dev/sda1 -> /dev/sda
func GetDeviceFromPartition(part InitService.MemoryPartition) InitService.MemoryDevice {
	partStr := string(part)

	// Standard Regex für die meisten Fälle
	re := regexp.MustCompile(`^(/dev/\D+?)(p?\d+)$`)
	matches := re.FindStringSubmatch(partStr)
	if len(matches) == 3 {
		return InitService.MemoryDevice(matches[1])
	}

	// Spezielle Behandlung für komplexere Device-Namen
	// z.B. /dev/mapper/vg-root → /dev/mapper/vg
	if strings.Contains(partStr, "/dev/mapper/") {
		re2 := regexp.MustCompile(`^(/dev/mapper/.+?)-(\w+)$`)
		matches2 := re2.FindStringSubmatch(partStr)
		if len(matches2) == 3 {
			return InitService.MemoryDevice(matches2[1])
		}
	}

	// Fallback: Original zurückgeben
	return InitService.MemoryDevice(part)
}

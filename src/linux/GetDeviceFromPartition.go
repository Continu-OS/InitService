package linux

import (
	"regexp"

	InitService "github.com/Continu-OS/InitService/src"
)

// Extrahiert das Basisspeichergerät aus der Partition, z.B.
// /dev/mmcblk0p1 -> /dev/mmcblk0
// /dev/sda1 -> /dev/sda
func GetDeviceFromPartition(part InitService.MemoryPartition) InitService.MemoryDevice {
	re := regexp.MustCompile(`^(/dev/\D+?)(p?\d+)$`)
	matches := re.FindStringSubmatch(string(part))
	if len(matches) == 3 {
		return InitService.MemoryDevice(matches[1])
	}
	return InitService.MemoryDevice(part)
}

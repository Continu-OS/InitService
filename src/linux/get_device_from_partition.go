package linux

import "regexp"

// Extrahiert das Basisspeichergerät aus der Partition, z.B.
// /dev/mmcblk0p1 -> /dev/mmcblk0
// /dev/sda1 -> /dev/sda
func GetDeviceFromPartition(part string) string {
	re := regexp.MustCompile(`^(/dev/\D+?)(p?\d+)$`)
	matches := re.FindStringSubmatch(part)
	if len(matches) == 3 {
		return matches[1]
	}
	return part
}

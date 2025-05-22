package linux

import "regexp"

// Extrahiert die Partitionsnummer aus einem Device-String,
// z.B. /dev/mmcblk0p1 -> "1", /dev/sda1 -> "1"
func ExtractPartitionNumber(dev string) string {
	re := regexp.MustCompile(`\d+$`)
	match := re.FindString(dev)
	return match
}

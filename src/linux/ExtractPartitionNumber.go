package linux

import (
	"regexp"

	InitService "github.com/Continu-OS/InitService/src"
)

// Extrahiert die Partitionsnummer aus einem Device-String,
// z.B. /dev/mmcblk0p1 -> "1", /dev/sda1 -> "1"
func ExtractPartitionNumber(parition InitService.MemoryPartition) string {
	re := regexp.MustCompile(`\d+$`)
	match := re.FindString(string(parition))
	return match
}

package linux

import (
	"log"
	"regexp"

	InitService "github.com/Continu-OS/syscored/src/syscored"
)

// Extrahiert die Partitionsnummer aus einem Device-String,
// z.B. /dev/mmcblk0p1 -> "1", /dev/sda1 -> "1"
func ExtractPartitionNumber(partition InitService.MemoryPartition) string {
	partStr := string(partition)

	// Standard: Letzte Ziffern
	re := regexp.MustCompile(`\d+$`)
	match := re.FindString(partStr)

	if match == "" {
		// Fallback für edge cases
		log.Printf("Warning: Could not extract partition number from %s", partStr)
		return ""
	}

	return match
}

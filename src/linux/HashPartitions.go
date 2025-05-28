package linux

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"strings"
)

// Diese Funktion erstellt aus der Partitionstabelle einen Hash
func HashPartitions(parts []PartitionInfo) string {
	var buf strings.Builder
	for _, p := range parts {
		// Füge relevante Felder zeilenweise hinzu (sortierbar, eindeutig)
		buf.WriteString(fmt.Sprintf("%s|%s|%s|%s\n", p.Name, p.DevicePath, p.Size, p.FsType))
	}
	sum := sha256.Sum256([]byte(buf.String()))
	return hex.EncodeToString(sum[:])
}

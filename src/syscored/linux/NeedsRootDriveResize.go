package linux

import (
	"fmt"

	InitService "github.com/Continu-OS/syscored/src/syscored"
)

// Prüft, ob ein Resize nötig ist (Partition kleiner als Gerät)
func NeedsRootDriveResize(device InitService.RootDevice, partition InitService.RootPartition) (bool, error) {
	devSize, err := GetSizeBytes(string(device))
	if err != nil {
		return false, fmt.Errorf("device size auslesen fehlgeschlagen: %v", err)
	}
	partSize, err := GetSizeBytes(string(partition))
	if err != nil {
		return false, fmt.Errorf("partition size auslesen fehlgeschlagen: %v", err)
	}

	return devSize > partSize, nil
}

package linux

import "fmt"

// Prüft, ob ein Resize nötig ist (Partition kleiner als Gerät)
func NeedsRootDriveResize(device, partition string) (bool, error) {
	devSize, err := GetSizeBytes(device)
	if err != nil {
		return false, fmt.Errorf("device size auslesen fehlgeschlagen: %v", err)
	}
	partSize, err := GetSizeBytes(partition)
	if err != nil {
		return false, fmt.Errorf("partition size auslesen fehlgeschlagen: %v", err)
	}

	return devSize > partSize, nil
}

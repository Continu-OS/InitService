package linux

import (
	"fmt"
	"os/exec"

	InitService "github.com/Continu-OS/InitService/src"
)

func ResizeRootDriveIfNeeded(device InitService.MemoryDevice, partition InitService.MemoryPartition) error {
	devSize, err := GetSizeBytes(string(device))
	if err != nil {
		return fmt.Errorf("device size auslesen fehlgeschlagen: %v", err)
	}
	partSize, err := GetSizeBytes(string(partition))
	if err != nil {
		return fmt.Errorf("partition size auslesen fehlgeschlagen: %v", err)
	}

	fmt.Printf("Gerätegröße: %d Bytes, Partitionsgröße: %d Bytes\n", devSize, partSize)

	if devSize > partSize {
		fmt.Println("Partition ist kleiner als Gerät, erweitere Partition...")

		partNum := ExtractPartitionNumber(partition)
		if partNum == "" {
			return fmt.Errorf("konnte Partitionsnummer nicht extrahieren aus %s", partition)
		}

		// Partition erweitern
		if err := exec.Command("growpart", string(device), partNum).Run(); err != nil {
			return fmt.Errorf("growpart fehlgeschlagen: %v", err)
		}

		// ext4-Dateisystem anpassen
		if err := exec.Command("resize2fs", string(partition)).Run(); err != nil {
			return fmt.Errorf("resize2fs fehlgeschlagen: %v", err)
		}

		fmt.Println("Partition und Dateisystem erfolgreich erweitert.")
	} else {
		fmt.Println("Partition ist bereits maximal, kein Resize nötig.")
	}
	return nil
}

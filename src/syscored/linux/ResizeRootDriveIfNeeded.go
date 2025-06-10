package linux

import (
	"fmt"
	"log"
	"os/exec"

	InitService "github.com/Continu-OS/syscored/src/syscored"
)

func ResizeRootDriveIfNeeded(device InitService.MemoryDevice, partition InitService.MemoryPartition) error {
	devSize, err := GetSizeBytes(string(device))
	if err != nil {
		return fmt.Errorf("failed to get device size for %s: %w", device, err)
	}

	partSize, err := GetSizeBytes(string(partition))
	if err != nil {
		return fmt.Errorf("failed to get partition size for %s: %w", partition, err)
	}

	log.Printf("Device size: %d bytes (%.2f GB), Partition size: %d bytes (%.2f GB)",
		devSize, float64(devSize)/(1024*1024*1024),
		partSize, float64(partSize)/(1024*1024*1024))

	// Threshold prüfen um unnötige Operationen zu vermeiden
	const minResizeThreshold = 100 * 1024 * 1024 // 100MB
	sizeDiff := devSize - partSize

	if sizeDiff <= minResizeThreshold {
		log.Printf("Size difference (%d bytes) below threshold (%d bytes), skipping resize",
			sizeDiff, minResizeThreshold)
		return nil
	}

	log.Printf("Partition is %d bytes smaller than device, expanding partition...", sizeDiff)

	// Partition Number extrahieren
	partNum := ExtractPartitionNumber(partition)
	if partNum == "" {
		return fmt.Errorf("could not extract partition number from %s", partition)
	}

	// Filesystem-Type erkennen vor dem Resize
	fsType, err := DetectFilesystemType(string(partition))
	if err != nil {
		log.Printf("Warning: Could not detect filesystem type: %v", err)
		fsType = "ext4" // Fallback
	}
	log.Printf("Detected filesystem type: %s", fsType)

	// 1. Partition erweitern mit growpart
	log.Printf("Expanding partition %s (partition %s)...", partition, partNum)
	growpartCmd := exec.Command("growpart", string(device), partNum)
	if output, err := growpartCmd.CombinedOutput(); err != nil {
		return fmt.Errorf("growpart failed: %w\nOutput: %s", err, string(output))
	}
	log.Println("Partition successfully expanded")

	// 2. Filesystem erweitern basierend auf Typ
	log.Printf("Resizing %s filesystem on %s...", fsType, partition)
	var resizeCmd *exec.Cmd

	switch fsType {
	case "ext2", "ext3", "ext4":
		resizeCmd = exec.Command("resize2fs", string(partition))
	case "xfs":
		// XFS braucht mount point, nicht device
		mountPoint, err := GetMountPoint(string(partition))
		if err != nil {
			return fmt.Errorf("failed to get mount point for XFS resize: %w", err)
		}
		resizeCmd = exec.Command("xfs_growfs", mountPoint)
	case "btrfs":
		mountPoint, err := GetMountPoint(string(partition))
		if err != nil {
			return fmt.Errorf("failed to get mount point for BTRFS resize: %w", err)
		}
		resizeCmd = exec.Command("btrfs", "filesystem", "resize", "max", mountPoint)
	default:
		return fmt.Errorf("unsupported filesystem type: %s", fsType)
	}

	if output, err := resizeCmd.CombinedOutput(); err != nil {
		return fmt.Errorf("filesystem resize failed: %w\nOutput: %s", err, string(output))
	}

	log.Printf("Filesystem successfully resized")

	// 3. Verification - neue Größen prüfen
	newPartSize, err := GetSizeBytes(string(partition))
	if err != nil {
		log.Printf("Warning: Could not verify new partition size: %v", err)
	} else {
		log.Printf("Resize completed: %d bytes (%.2f GB) -> %d bytes (%.2f GB)",
			partSize, float64(partSize)/(1024*1024*1024),
			newPartSize, float64(newPartSize)/(1024*1024*1024))
	}

	return nil
}

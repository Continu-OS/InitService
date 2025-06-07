package linux

import (
	"fmt"
	"os/exec"
	"strings"
)

// Hilfsfunktion für Filesystem-Type Detection
func DetectFilesystemType(device string) (string, error) {
	cmd := exec.Command("lsblk", "-no", "FSTYPE", device)
	output, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("lsblk failed: %w", err)
	}

	fsType := strings.TrimSpace(string(output))
	if fsType == "" {
		return "", fmt.Errorf("no filesystem type detected")
	}

	return fsType, nil
}

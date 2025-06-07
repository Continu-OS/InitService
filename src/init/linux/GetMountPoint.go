package linux

import (
	"fmt"
	"os/exec"
	"strings"
)

// Hilfsfunktion für Mount Point Detection
func GetMountPoint(device string) (string, error) {
	cmd := exec.Command("lsblk", "-no", "MOUNTPOINT", device)
	output, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("lsblk failed: %w", err)
	}

	mountPoint := strings.TrimSpace(string(output))
	if mountPoint == "" {
		return "", fmt.Errorf("device not mounted")
	}

	return mountPoint, nil
}

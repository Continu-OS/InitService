package linux

import (
	"fmt"
	"os/exec"
	"strconv"
	"strings"
)

// Liest die Größe in Bytes von Gerät oder Partition mit blockdev --getsize64
func GetSizeBytes(device string) (int64, error) {
	cmd := exec.Command("blockdev", "--getsize64", device)
	out, err := cmd.Output()
	if err != nil {
		return 0, fmt.Errorf("failed to get size for device %s: %w", device, err)
	}

	sizeStr := strings.TrimSpace(string(out))
	if sizeStr == "" {
		return 0, fmt.Errorf("empty size output for device %s", device)
	}

	size, err := strconv.ParseInt(sizeStr, 10, 64)
	if err != nil {
		return 0, fmt.Errorf("failed to parse size '%s' for device %s: %w", sizeStr, device, err)
	}

	return size, nil
}

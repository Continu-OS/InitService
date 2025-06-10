package linux

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	InitService "github.com/Continu-OS/syscored/src/syscored"
)

func GetRootDevicePartition() (InitService.RootPartition, error) {
	f, err := os.Open("/proc/mounts")
	if err != nil {
		return "", fmt.Errorf("failed to open /proc/mounts: %w", err)
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue // Skip empty lines and comments
		}

		fields := strings.Fields(line)
		if len(fields) >= 2 && fields[1] == "/" {
			device := fields[0]
			// Validate device path
			if !strings.HasPrefix(device, "/dev/") {
				continue // Skip non-device mounts (tmpfs, proc, etc.)
			}
			return InitService.RootPartition(device), nil
		}
	}

	if err := scanner.Err(); err != nil {
		return "", fmt.Errorf("error reading /proc/mounts: %w", err)
	}

	return "", fmt.Errorf("root device not found in /proc/mounts")
}

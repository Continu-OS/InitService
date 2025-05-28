package linux

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	InitService "github.com/Continu-OS/InitService/src"
)

func GetRootDevicePartition() (InitService.RootPartition, error) {
	f, err := os.Open("/proc/mounts")
	if err != nil {
		return "", err
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := scanner.Text()
		fields := strings.Fields(line)
		if len(fields) >= 2 && fields[1] == "/" {
			return InitService.RootPartition(fields[0]), nil
		}
	}
	if err := scanner.Err(); err != nil {
		return "", err
	}
	return "", fmt.Errorf("root device not found")
}

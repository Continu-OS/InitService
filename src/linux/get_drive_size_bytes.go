package linux

import (
	"os/exec"
	"strconv"
	"strings"
)

// Liest die Größe in Bytes von Gerät oder Partition mit blockdev --getsize64
func GetSizeBytes(device string) (int64, error) {
	out, err := exec.Command("blockdev", "--getsize64", device).Output()
	if err != nil {
		return 0, err
	}
	sizeStr := strings.TrimSpace(string(out))
	return strconv.ParseInt(sizeStr, 10, 64)
}

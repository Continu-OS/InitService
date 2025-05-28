package linux

import (
	"os"
	"strings"

	InitService "github.com/Continu-OS/InitService/src"
)

func VerifyRootSystemDriveSignature(rootDev InitService.RootPartition, device InitService.MemoryDevice) (bool, error) {
	parts, err := ListPartitions(device)
	if err != nil {
		return false, err
	}

	partsHash := HashPartitions(parts)

	sig, err := os.ReadFile(InitService.HostSystemRootDriveValidationFilePath)
	if err != nil {
		return false, err
	}

	expected := strings.TrimSpace(string(sig))

	if partsHash == expected {
		return true, nil
	}

	return false, nil
}

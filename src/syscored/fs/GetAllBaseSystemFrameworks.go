package fs

import (
	"fmt"
	"os"

	InitService "github.com/Continu-OS/syscored/src/syscored"
)

func GetAllBaseSystemFrameworks() ([]*InitService.BaseSystemFramework, error) {
	entries, err := os.ReadDir(InitService.HostBaseSystemServicesDirPath)
	if err != nil {
		return nil, err
	}
	for _, entry := range entries {
		if entry.IsDir() {
			fmt.Println("Ordner gefunden:", entry.Name())
		}
	}
	return nil, nil
}

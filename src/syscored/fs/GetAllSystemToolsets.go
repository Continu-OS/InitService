package fs

import (
	"fmt"
	"os"

	InitService "github.com/Continu-OS/syscored/src/syscored"
)

func GetAllSystemToolsets() ([]*InitService.SystemToolset, error) {
	entries, err := os.ReadDir(InitService.HostSystemToolsetsDirPath)
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

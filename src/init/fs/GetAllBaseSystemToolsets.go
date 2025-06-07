package fs

import (
	"fmt"
	"os"

	InitService "github.com/Continu-OS/InitService/src/init"
)

func GetAllBaseSystemToolsets() ([]*InitService.BaseSystemToolset, error) {
	entries, err := os.ReadDir(InitService.HostBasesystemToolsetsDirPath)
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

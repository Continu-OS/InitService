package firmware

import (
	"log"
	"os"
	"strings"

	InitService "github.com/Continu-OS/InitService/src"
)

// GetAllBootloaderParameters liest alle Boot-Parameter aus /proc/cmdline aus
func GetAllBootloaderParameters() (InitService.BootloaderStartParameters, error) {
	data, err := os.ReadFile("/proc/cmdline")
	if err != nil {
		return nil, err
	}

	log.Println("Extracting Bootloader Arguments")

	cmdline := string(data)
	params := strings.Fields(cmdline)

	result := make(map[string]string)
	for _, param := range params {
		if strings.Contains(param, "=") {
			parts := strings.SplitN(param, "=", 2)
			key := parts[0]
			value := parts[1]
			result[key] = value
		} else {
			result[param] = ""
		}
	}

	return result, nil
}

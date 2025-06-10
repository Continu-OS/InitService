package firmware

import (
	"fmt"
	"log"
	"os"
	"strings"

	InitService "github.com/Continu-OS/syscored/src/syscored"
)

// GetAllBootloaderParameters liest alle Boot-Parameter aus /proc/cmdline aus
// Kompatibel mit allen Linux-Bootloadern (GRUB, U-Boot, systemd-boot, etc.)
func GetAllBootloaderParameters() (InitService.BootloaderStartParameters, error) {
	data, err := os.ReadFile("/proc/cmdline")
	if err != nil {
		return nil, fmt.Errorf("failed to read /proc/cmdline: %w", err)
	}

	log.Println("Extracting Bootloader Arguments")

	// Whitespace normalisieren und trimmen
	cmdline := strings.TrimSpace(string(data))
	if cmdline == "" {
		log.Println("Warning: Empty cmdline detected")
		return make(map[string]string), nil
	}

	params := strings.Fields(cmdline)
	result := make(map[string]string)

	// Für Parameter die mehrfach vorkommen können
	multiParams := make(map[string][]string)

	for _, param := range params {
		param = strings.TrimSpace(param)
		if param == "" {
			continue
		}

		if strings.Contains(param, "=") {
			parts := strings.SplitN(param, "=", 2)
			key := strings.TrimSpace(parts[0])
			value := strings.TrimSpace(parts[1])

			// Quotes entfernen falls vorhanden (U-Boot Kompatibilität)
			if len(value) >= 2 &&
				((value[0] == '"' && value[len(value)-1] == '"') ||
					(value[0] == '\'' && value[len(value)-1] == '\'')) {
				value = value[1 : len(value)-1]
			}

			// Behandlung mehrfacher Parameter
			if existing, exists := result[key]; exists {
				if multiValues, isMulti := multiParams[key]; isMulti {
					multiParams[key] = append(multiValues, value)
				} else {
					multiParams[key] = []string{existing, value}
				}
				result[key] = strings.Join(multiParams[key], ",")
			} else {
				result[key] = value
			}
		} else {
			// Flag ohne Wert
			result[param] = ""
		}
	}

	// Debug-Logging falls verfügbar
	log.Printf("Parsed %d bootloader parameters from cmdline", len(result))

	// Optional: Detailliertes Debug-Logging
	if os.Getenv("DEBUG_BOOTLOADER_PARAMS") == "1" {
		log.Printf("Raw cmdline: %s", cmdline)
		for k, v := range result {
			if v == "" {
				log.Printf("  Flag: %s", k)
			} else {
				log.Printf("  Param: %s=%s", k, v)
			}
		}
	}

	return result, nil
}

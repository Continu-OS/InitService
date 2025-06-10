package firmware

import (
	"fmt"
	"os"
)

func CheckEFISecurebootIsEnable() {
	path := "/sys/firmware/efi/efivars/SecureBoot-8be4df61-93ca-11d2-aa0d-00e098032b8c"

	data, err := os.ReadFile(path)
	if err != nil {
		fmt.Println("Secure Boot-Status nicht lesbar oder kein UEFI-Modus:", err)
		return
	}

	// Die ersten 4 Bytes sind Attribute, der fünfte ist der Wert
	if len(data) < 5 {
		fmt.Println("Unerwarteter Dateninhalt.")
		return
	}

	if data[4] == 1 {
		fmt.Println("Secure Boot ist AKTIV.")
	} else if data[4] == 0 {
		fmt.Println("Secure Boot ist INAKTIV.")
	} else {
		fmt.Printf("Unerwarteter Wert: %d\n", data[4])
	}
}

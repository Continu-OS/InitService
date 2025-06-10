package linux

import (
	"log"
	"os"

	"github.com/Continu-OS/syscored/src/syscored/linux/cgroups"
)

func RebootSystemNow() error {
	if cgroups.IsRunningInDocker() {
		log.Println("Docker detected – skipping system reboot (exiting process instead)")
		os.Exit(0)
		return nil
	}
	return nil
}

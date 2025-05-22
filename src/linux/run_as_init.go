package linux

import (
	"os"
	"os/user"
)

func RunAsInitProcess() (bool, error) {
	// Prüfen ob PID 1
	pid := os.Getpid()
	if pid != 1 {
		return false, nil
	}

	// Prüfen ob als Root
	u, err := user.Current()
	if err != nil {
		return false, err
	}

	if u.Uid != "0" {
		return false, nil
	}

	// Alternative: UID direkt über syscall
	uid := os.Geteuid()
	if uid != 0 {
		return false, nil
	}

	return true, nil
}

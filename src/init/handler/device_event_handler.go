package handler

import (
	"log"
	"strings"

	"github.com/Continu-OS/InitService/src/init/linux/cgroups"
	"github.com/mdlayher/netlink"
	"golang.org/x/sys/unix"
)

// StartDeviceEventListener listens for device events (e.g., block, USB, net)
// and logs basic uevent information.
func StartDeviceEventHandler() error {
	if cgroups.IsRunningInDocker() {
		log.Println("Docker detected – skipping device event handler startup")
		return nil
	}

	conn, err := netlink.Dial(unix.NETLINK_KOBJECT_UEVENT, nil)
	if err != nil {
		return err
	}

	go func() {
		defer conn.Close()

		for {
			msgs, err := conn.Receive()
			if err != nil {
				log.Println("DeviceEventListener error:", err)
				continue
			}

			for _, m := range msgs {
				lines := parseUeventPayload(m.Data)
				handleUevent(lines)
			}
		}
	}()
	return nil
}

// parseUeventPayload splits the raw netlink data into separate key=value strings.
func parseUeventPayload(data []byte) []string {
	var result []string
	start := 0
	for i, b := range data {
		if b == 0 {
			result = append(result, string(data[start:i]))
			start = i + 1
		}
	}
	return result
}

// handleUevent logs selected device event fields.
func handleUevent(lines []string) {
	log.Println("Netlink uevent received:")
	for _, line := range lines {
		if strings.HasPrefix(line, "ACTION=") ||
			strings.HasPrefix(line, "DEVNAME=") ||
			strings.HasPrefix(line, "SUBSYSTEM=") {
			log.Println("  ", line)
		}
	}
}

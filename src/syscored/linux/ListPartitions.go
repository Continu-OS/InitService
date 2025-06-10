package linux

import (
	"bytes"
	"os/exec"
	"strings"

	InitService "github.com/Continu-OS/syscored/src/syscored"
)

func ListPartitions(device InitService.MemoryDevice) ([]PartitionInfo, error) {
	cmd := exec.Command("lsblk", "-o", "NAME,SIZE,FSTYPE,TYPE", "-nr", string(device))
	out, err := cmd.Output()
	if err != nil {
		return nil, err
	}

	var partitions []PartitionInfo
	lines := bytes.Split(out, []byte("\n"))
	for _, line := range lines {
		fields := strings.Fields(string(line))
		if len(fields) != 4 {
			continue
		}
		name, size, fstype, typ := fields[0], fields[1], fields[2], fields[3]
		if typ == "part" {
			partitions = append(partitions, PartitionInfo{
				Name:       name,
				DevicePath: "/dev/" + name,
				Size:       size,
				FsType:     fstype,
			})
		}
	}
	return partitions, nil
}

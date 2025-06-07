package linux

// PartitionInfo enthält Informationen zu einer Partition
type PartitionInfo struct {
	Name       string // z.B. "sda1"
	DevicePath string // z.B. "/dev/sda1"
	Size       string // z.B. "100G"
	FsType     string // z.B. "ext4"
}

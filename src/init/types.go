package InitService

type FQP string
type ToolsetBaseLevel uint8
type BootloaderStartParameters map[string]string
type MemoryDevice string
type MemoryPartition string
type RootDevice MemoryDevice
type RootPartition MemoryPartition

type ToolsetModule struct {
	Path      FQP
	Name      string
	BaseLevel ToolsetBaseLevel
}

type HostInitConfig struct {
}

type BaseSystemService struct {
}

type BaseSystemFramework struct {
}

type BaseSystemToolset struct {
}

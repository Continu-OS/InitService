package InitService

type FQP string
type ToolsetBaseLevel uint8

type ToolsetModule struct {
	Path      FQP
	Name      string
	BaseLevel ToolsetBaseLevel
}

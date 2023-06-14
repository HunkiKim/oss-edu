package pkg

type IoType int

const (
	Cli IoType = iota
	File
	Json
)

type User struct {
	Name          string
	NumberOfTurns int
}

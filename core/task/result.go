package task

type Status int

const (
	OK Status = iota
	FAILURE
)

type Result struct {
	TaskID string
	Data   []byte
	Status Status
}

package executor

type Status string

const (
	OK     Status = "ok"
	Failed        = "failed"
)

type Result struct {
	TaskID string
	Status Status
	Data   []byte
}

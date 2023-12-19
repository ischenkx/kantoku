package alpha

type Status int

const (
	OK      Status = 0
	FAILURE        = 1
)

type Result struct {
	Data   []byte
	Status Status
}

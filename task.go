package kantoku

import "context"

type Task struct {
	ID_          string
	Type_        string
	Argument_    []byte
	Dependencies []string
}

func (t Task) ID(ctx context.Context) string {
	return t.ID_
}

func (t Task) Type(ctx context.Context) string {
	return t.Type_
}

func (t Task) Argument(ctx context.Context) []byte {
	return t.Argument_
}

const (
	ArgumentTypeCell = iota
	ArgumentTypeConstant
	ArgumentTypeTaskResult
)

// Argument probably should be elsewhere
type Argument struct {
	Type_ int
	Value string
}

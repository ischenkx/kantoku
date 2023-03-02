package kantoku

import (
	"context"
	"kantoku/core/l1"
)

type Task struct {
	Spec         l1.Task
	Dependencies []string
}

func (t Task) L1(ctx context.Context) (l1.Task, error) {
	return t.Spec, nil
}

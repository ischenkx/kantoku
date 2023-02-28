package kantoku

import (
	"context"
	"kantoku/core/l1"
	"kantoku/framework/depot"
)

type Task struct {
	Spec         l1.Task
	Dependencies []string
}

func (t Task) AsL1(ctx context.Context) (l1.Task, error) {
	return t.Spec, nil
}

func (t Task) AsDepot() depot.Task {
	return depot.Task{
		ID:           t.Spec.ID,
		Dependencies: t.Dependencies,
	}
}

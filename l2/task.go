package l2

import (
	"context"
	"kantoku/l1"
)

type Task interface {
	ToL1(ctx context.Context) (l1.Task, error)
}

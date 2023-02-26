package l2

import (
	"context"
	"kantoku/core/l1"
)

type Task interface {
	AsL1(ctx context.Context) (l1.Task, error)
}

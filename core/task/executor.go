package task

import "context"

type Executor[Task any] interface {
	Execute(ctx context.Context, task Task) (Result, error)
}

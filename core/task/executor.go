package task

import "context"

type Executor[Task AbstractTask] interface {
	Execute(ctx context.Context, task Task) (Result, error)
}

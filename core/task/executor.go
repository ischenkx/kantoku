package task

import "context"

type Executor[ArgumentType any] interface {
	Execute(ctx context.Context, task AbstractTask) (Result, error)
}

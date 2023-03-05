package task

import "context"

type Executor[InputType AbstractTask] interface {
	Execute(ctx context.Context, task InputType) (Result, error)
}

package executor

import "context"

type Runner[Task, Output any] interface {
	Run(context.Context, Task) (Output, error)
}

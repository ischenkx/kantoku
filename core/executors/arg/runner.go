package arg

import (
	"context"
	"kantoku/core/task"
)

type Runner interface {
	Run(ctx context.Context, function string, args []Arg) (task.Result, error)
}

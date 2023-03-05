package simple

import (
	"context"
	"errors"
	"kantoku/common/util"
	"kantoku/core/task"
)

type Executor[Input task.AbstractTask] map[string]func([]byte) ([]byte, error)

func (e Executor[Input]) Execute(ctx context.Context, input Input) (task.Result, error) {
	executor, ok := e[input.Type(ctx)]
	if !ok {
		return util.Default[task.Result](), errors.New("no task for a given type: " + input.Type(ctx))
	}

	data, err := executor(input.Argument(ctx))
	if err != nil {
		return task.Result{
			TaskID: input.ID(ctx),
			Data:   []byte(err.Error()),
			Status: task.FAILURE,
		}, nil
	}

	return task.Result{
		TaskID: input.ID(ctx),
		Data:   data,
		Status: task.OK,
	}, nil
}

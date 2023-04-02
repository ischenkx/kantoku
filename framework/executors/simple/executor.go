package simple

import (
	"context"
	"errors"
	"kantoku"
	"kantoku/common/util"
	"kantoku/core/task"
)

type Executor map[string]func(any) ([]byte, error)

func (e Executor) Execute(ctx context.Context, input kantoku.StoredTask) (task.Result, error) {
	executor, ok := e[input.Type]
	if !ok {
		return util.Default[task.Result](), errors.New("no executor for a given type: " + input.Type)
	}

	data, err := executor(input.Data)
	if err != nil {
		return task.Result{
			TaskID: input.Id,
			Data:   []byte(err.Error()),
			Status: task.FAILURE,
		}, nil
	}

	return task.Result{
		TaskID: input.Id,
		Data:   data,
		Status: task.OK,
	}, nil
}

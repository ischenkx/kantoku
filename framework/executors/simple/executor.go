package simple

import (
	"context"
	"errors"
	"kantoku"
	"kantoku/common/util"
	"kantoku/core/task"
)

type Executor map[string]func(ctx context.Context, task *kantoku.View) ([]byte, error)

func (e Executor) Execute(ctx context.Context, view *kantoku.View) (task.Result, error) {
	storedTask, err := view.Stored(ctx)
	if err != nil {
		return task.Result{}, err
	}

	executor, ok := e[storedTask.Type]
	if !ok {
		return util.Default[task.Result](), errors.New("no executor for a given type: " + storedTask.Type)
	}

	data, err := executor(ctx, view)
	if err != nil {
		return task.Result{
			TaskID: storedTask.Id,
			Data:   []byte(err.Error()),
			Status: task.FAILURE,
		}, nil
	}

	return task.Result{
		TaskID: storedTask.Id,
		Data:   data,
		Status: task.OK,
	}, nil
}

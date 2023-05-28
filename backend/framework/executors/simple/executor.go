package simple

import (
	"context"
	"errors"
	"kantoku"
	"kantoku/common/util"
	"kantoku/platform"
)

type Executor map[string]func(ctx context.Context, task *kantoku.View) ([]byte, error)

func (e Executor) Execute(ctx context.Context, view *kantoku.View) (platform.Result, error) {
	TaskInstance, err := view.Stored(ctx)
	if err != nil {
		return platform.Result{}, err
	}

	executor, ok := e[TaskInstance.Type]
	if !ok {
		return util.Default[platform.Result](), errors.New("no executor for a given type: " + TaskInstance.Type)
	}

	data, err := executor(ctx, view)
	if err != nil {
		return platform.Result{
			TaskID: TaskInstance.Id,
			Data:   []byte(err.Error()),
			Status: platform.FAILURE,
		}, nil
	}

	return platform.Result{
		TaskID: TaskInstance.Id,
		Data:   data,
		Status: platform.OK,
	}, nil
}

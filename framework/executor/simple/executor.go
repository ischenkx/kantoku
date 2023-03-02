package simple

import (
	"errors"
	"kantoku/common/util"
	"kantoku/core/l1"
)

type Executor map[string]func([]byte) ([]byte, error)

func (e Executor) Execute(task l1.Task) (l1.Result, error) {
	executor, ok := e[task.Type]
	if !ok {
		return util.Default[l1.Result](), errors.New("no executor for a given type: " + task.Type)
	}

	data, err := executor(task.Argument)
	if err != nil {
		return l1.Result{
			TaskID: task.ID,
			Data:   []byte(err.Error()),
			Status: l1.FAILURE,
		}, nil
	}

	return l1.Result{
		TaskID: task.ID,
		Data:   data,
		Status: l1.OK,
	}, nil
}

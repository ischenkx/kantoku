package converters

import (
	"github.com/ischenkx/kantoku/pkg/core/task"
	"github.com/ischenkx/kantoku/pkg/lib/connector/api/http/oas"
)

func TaskToDto(t task.Task) oas.Task {
	return oas.Task{
		Id:      t.ID,
		Info:    t.Info,
		Inputs:  t.Inputs,
		Outputs: t.Outputs,
	}
}

package kantokuhttp

import (
	"github.com/ischenkx/kantoku/pkg/core"
	"github.com/ischenkx/kantoku/pkg/lib/gateway/api/kantokuhttp/oas"
)

func TaskToDto(t core.Task) oas.Task {
	return oas.Task{
		Id:      t.ID,
		Info:    t.Info,
		Inputs:  t.Inputs,
		Outputs: t.Outputs,
	}
}

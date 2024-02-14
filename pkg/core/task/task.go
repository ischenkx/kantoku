package task

import (
	"github.com/ischenkx/kantoku/pkg/common/data/record"
	"github.com/ischenkx/kantoku/pkg/core/resource"
)

type Task struct {
	Inputs  []resource.ID
	Outputs []resource.ID
	ID      string
	Info    record.R
}

func New(options ...Option) Task {
	t := Task{Info: record.R{}}

	for _, option := range options {
		option(&t)
	}

	return t
}

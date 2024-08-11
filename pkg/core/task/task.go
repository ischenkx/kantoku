package task

import (
	"github.com/ischenkx/kantoku/pkg/common/data/storage"
	"github.com/ischenkx/kantoku/pkg/core/resource"
)

type Task struct {
	Inputs  []resource.ID
	Outputs []resource.ID
	ID      string
	Info    map[string]any
}

func New(options ...Option) Task {
	t := Task{Info: map[string]any{}}

	for _, option := range options {
		option(&t)
	}

	return t
}

func (task Task) AsDoc() storage.Document {
	return map[string]any{
		"id":      task.ID,
		"inputs":  task.Inputs,
		"outputs": task.Outputs,
		"info":    task.Info,
	}
}

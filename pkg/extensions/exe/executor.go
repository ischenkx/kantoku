package exe

import (
	"context"
	"github.com/ischenkx/kantoku/pkg/system"
	"github.com/ischenkx/kantoku/pkg/system/kernel/task"
)

type Executor struct {
	fn func(*Context) error
}

func New(fn func(*Context) error) *Executor {
	return &Executor{fn: fn}
}

func (executor *Executor) Execute(ctx context.Context, sys system.AbstractSystem, task task.Task) error {
	return executor.fn(&Context{
		system: sys,
		task:   task,
		ctx:    ctx,
	})
}

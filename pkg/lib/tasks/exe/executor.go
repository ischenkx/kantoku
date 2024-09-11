package exe

import (
	"context"
	"github.com/ischenkx/kantoku/pkg/core"
)

type Executor struct {
	fn func(*Context) error
}

func New(fn func(*Context) error) *Executor {
	return &Executor{fn: fn}
}

func (executor *Executor) Execute(ctx context.Context, sys core.AbstractSystem, task core.Task) error {
	return executor.fn(&Context{
		system: sys,
		task:   task,
		ctx:    ctx,
	})
}

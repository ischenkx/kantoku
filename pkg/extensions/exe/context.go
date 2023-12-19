package exe

import (
	"context"
	"github.com/ischenkx/kantoku/pkg/system"
	"github.com/ischenkx/kantoku/pkg/system/kernel/task"
	"time"
)

type Context struct {
	system system.AbstractSystem
	task   task.Task
	ctx    context.Context
}

func (context *Context) System() system.AbstractSystem {
	return context.system
}

func (context *Context) Task() task.Task {
	return context.task
}

// ########## context.Context ##########

func (context *Context) Deadline() (deadline time.Time, ok bool) {
	return context.ctx.Deadline()
}

func (context *Context) Done() <-chan struct{} {
	return context.ctx.Done()
}

func (context *Context) Err() error {
	return context.ctx.Err()
}

func (context *Context) Value(key any) any {
	return context.ctx.Value(key)
}

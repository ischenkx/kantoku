package exe

import (
	"context"
	"github.com/ischenkx/kantoku/pkg/core"
	"time"
)

type Context struct {
	system core.AbstractSystem
	task   core.Task
	ctx    context.Context
}

func (context *Context) System() core.AbstractSystem {
	return context.system
}

func (context *Context) Task() core.Task {
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

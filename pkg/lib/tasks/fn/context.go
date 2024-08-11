package fn

import (
	"context"
	"github.com/ischenkx/kantoku/pkg/lib/tasks/fn_d/future"
)

type ScheduledTask struct {
	Type    string
	Inputs  []future.AbstractFuture
	Outputs []future.AbstractFuture
}

type Context struct {
	context.Context

	futures map[int32]future.AbstractFuture

	scheduledTasks []ScheduledTask
}

func NewContext(ctx context.Context) *Context {
	return &Context{
		Context: ctx,
	}
}

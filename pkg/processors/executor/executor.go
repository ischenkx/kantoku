package executor

import (
	"context"
	"github.com/ischenkx/kantoku/pkg/system"
	"github.com/ischenkx/kantoku/pkg/system/kernel/task"
)

type Executor interface {
	Execute(ctx context.Context, sys system.AbstractSystem, task task.Task) error
}

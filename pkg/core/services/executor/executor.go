package executor

import (
	"context"
	"github.com/ischenkx/kantoku/pkg/core/system"
	"github.com/ischenkx/kantoku/pkg/core/task"
)

type Executor interface {
	Execute(ctx context.Context, sys system.AbstractSystem, task task.Task) error
}

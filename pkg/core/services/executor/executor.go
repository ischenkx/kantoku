package executor

import (
	"context"
	"github.com/ischenkx/kantoku/pkg/core"
)

type Executor interface {
	Execute(ctx context.Context, sys core.AbstractSystem, task core.Task) error
}

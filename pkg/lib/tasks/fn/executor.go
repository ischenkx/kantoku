package fn

import (
	"context"
	"fmt"
	"github.com/ischenkx/kantoku/pkg/core/services/executor"
	"github.com/ischenkx/kantoku/pkg/core/system"
	"github.com/ischenkx/kantoku/pkg/core/task"
)

var _ executor.Executor = (*Executor)(nil)

type Executor[F AbstractFunction[I, O], I, O any] struct {
	function  F
	scheduler *Scheduler
}

func NewExecutor[F AbstractFunction[I, O], I, O any](f F) (*Executor[F, I, O], error) {
	return &Executor[F, I, O]{
		function:  f,
		scheduler: &Scheduler{},
	}, nil
}

func (e *Executor[T, I, O]) Execute(ctx context.Context, sys system.AbstractSystem, task task.Task) error {
	// Allocate a new context
	fctx := NewContext(ctx)

	// Bind task inputs to the "I" type

	var input I

	output, err := e.function.Call(fctx, input)
	if err != nil {
		return fmt.Errorf("failed to execute function: %w", err)
	}

	// Save output to context

	// Allocate Resources

	// Spawn Tasks from context

	return nil
}

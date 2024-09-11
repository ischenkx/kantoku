package recursive

import (
	"fmt"
	"github.com/ischenkx/kantoku/pkg/lib/tasks/fn"
	"github.com/ischenkx/kantoku/pkg/lib/tasks/fn/future"
)

type (
	BInput struct {
		Length future.Future[int]
	}

	BOutput struct {
		Calc future.Future[int]
	}

	B struct {
		fn.Function[*B, BInput, BOutput]
	}
)

var (
	_ fn.AbstractFunction[BInput, BOutput] = (*B)(nil)
)

func (a B) Call(ctx *fn.Context, input BInput) (output BOutput, err error) {
	length := input.Length.Value()

	if length <= 0 {
		return BOutput{Calc: future.FromValue(-1)}, nil
	}

	length--

	_, err = fn.Sched[A](ctx, AInput{
		Length: future.FromValue(length),
	})
	if err != nil {
		return output, fmt.Errorf("failed to schedule A: %w", err)
	}

	_, err = fn.Sched[A](ctx, AInput{
		Length: future.FromValue(length / 2),
	})
	if err != nil {
		return output, fmt.Errorf("failed to schedule A: %w", err)
	}

	return BOutput{Calc: future.FromValue(1)}, nil
}

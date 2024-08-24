package recursive

import (
	"fmt"
	"github.com/ischenkx/kantoku/pkg/lib/tasks/fn"
	"github.com/ischenkx/kantoku/pkg/lib/tasks/fn/future"
)

type (
	AInput struct {
		Length future.Future[int]
	}

	AOutput struct {
		Calc future.Future[int]
	}

	A struct {
		fn.Function[A, AInput, AOutput]
	}
)

var (
	_ fn.AbstractFunction[AInput, AOutput] = (*A)(nil)
)

func (a A) Call(ctx *fn.Context, input AInput) (output AOutput, err error) {
	length := input.Length.Value()

	if length <= 0 {
		return AOutput{Calc: future.FromValue(-1)}, nil
	}

	length--

	o1, err := fn.Sched[B](ctx, BInput{
		Length: future.FromValue(length),
	})
	if err != nil {
		return output, fmt.Errorf("failed to schedule B: %w", err)
	}

	_, err = fn.Sched[B](ctx, BInput{
		Length: o1.Calc,
	})
	if err != nil {
		return output, fmt.Errorf("failed to schedule B: %w", err)
	}

	return AOutput{Calc: future.FromValue(1)}, nil
}

package main

import (
	"context"
	"github.com/ischenkx/kantoku/cmd/testing/stand"
	"github.com/ischenkx/kantoku/pkg/lib/tasks/functional"
	"github.com/ischenkx/kantoku/pkg/lib/tasks/future"
)

func main() {
	ctx := context.Background()
	sys := NewSystem(ctx)

	_, err := functional.Do(context.Background(), sys, func(ctx *functional.Context) error {
		for i := 0; i < 1; i++ {
			functional.Execute[stand.SumTask, stand.SumInput, stand.MathOutput](ctx, stand.SumTask{},
				stand.SumInput{Args: future.FromValue([]int{1, 2, 3, 4, 5, 6, 7, 8, 9})},
			)
		}

		return nil
	})
	if err != nil {
		panic(err)
	}
}

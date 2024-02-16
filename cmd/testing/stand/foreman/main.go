package main

import (
	"context"
	"github.com/ischenkx/kantoku/cmd/testing/stand"
	"github.com/ischenkx/kantoku/cmd/testing/stand/common"
	"github.com/ischenkx/kantoku/pkg/lib/tasks/functional"
	"github.com/ischenkx/kantoku/pkg/lib/tasks/future"
)

func main() {
	ctx := context.Background()
	sys := stand.NewSystem(ctx)

	err := functional.SchedulingContext(context.Background(), sys, func(ctx *functional.Context) error {
		for i := 0; i < 1; i++ {
			functional.Execute[common.SumTask, common.SumInput, common.MathOutput](ctx, common.SumTask{},
				common.SumInput{Args: future.FromValue([]int{1, 2, 3, 4, 5, 6, 7, 8, 9})},
			)
		}

		return nil
	})
	if err != nil {
		panic(err)
	}
}

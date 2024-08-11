package main

import (
	"context"
	"github.com/ischenkx/kantoku/cmd/stand/utils"
	main2 "github.com/ischenkx/kantoku/local/deprecated/math_executor"
	"github.com/ischenkx/kantoku/pkg/lib/platform"
	"github.com/ischenkx/kantoku/pkg/lib/tasks/fn_d"
	"github.com/ischenkx/kantoku/pkg/lib/tasks/fn_d/future"
	"log"
	"os"
)

func main() {
	ctx := context.Background()
	cfg := utils.LoadConfig()
	logger := utils.GetLogger(os.Stdout, "foreman")

	sys, err := platform.BuildSystem(ctx, logger, cfg.Core.System)
	if err != nil {
		log.Fatal("failed to build system: ", err)
	}

	err = fn_d.SchedulingContext(context.Background(), sys, func(ctx *fn_d.Context) error {
		for i := 0; i < 10; i++ {
			fn_d.Execute[main2.SumTask, main2.SumInput, main2.MathOutput](ctx, main2.SumTask{},
				main2.SumInput{Args: future.FromValue([]int{1, 2, 3, 4, 5, 6, 7, 8, 9})},
			)
		}

		return nil
	})
	if err != nil {
		panic(err)
	}
}

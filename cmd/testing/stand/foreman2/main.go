package main

import (
	"context"
	"github.com/ischenkx/kantoku/cmd/testing/stand/common"
	"github.com/ischenkx/kantoku/pkg/common/data/codec"
	"github.com/ischenkx/kantoku/pkg/lib/tasks/functional"
	"github.com/ischenkx/kantoku/pkg/lib/tasks/future"
	"log"
)

func main() {
	common.InitLogger()
	ctx := context.Background()
	sys := common.NewSystem(ctx, "foreman-func")

	res, err := sys.Resources().Load(context.Background(), "c3469a9766e94e4db7e3d6fa8fa86734")
	if err != nil {
		panic(err)
	}
	decoded, err := codec.JSON[int]().Decode(res[0].Data)
	if err != nil {
		panic(err)
	}
	log.Println(decoded)
	err = functional.SchedulingContext(context.Background(), sys, func(ctx *functional.Context) error {
		for i := 0; i < 1000; i++ {
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

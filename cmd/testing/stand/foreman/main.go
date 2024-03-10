package main

import (
	"context"
	"fmt"
	"github.com/ischenkx/kantoku/cmd/testing/stand"
	"github.com/ischenkx/kantoku/pkg/common/data/record"
	"github.com/ischenkx/kantoku/pkg/common/data/record/ops"
	recutil "github.com/ischenkx/kantoku/pkg/common/data/record/util"
	"github.com/ischenkx/kantoku/pkg/core/resource"
	"github.com/ischenkx/kantoku/pkg/core/task"
	"github.com/ischenkx/kantoku/pkg/lib/tasks/functional"
	"github.com/ischenkx/kantoku/pkg/lib/tasks/future"
	"github.com/samber/lo"
	"math/rand"
	"time"
)

func main() {
	ctx := context.Background()
	sys := NewSystem(ctx)

	var result int
	taskCtx, err := functional.Do(context.Background(), sys, func(ctx *functional.Context) error {
		var previousFuture future.Future[int]

		for i := 0; i < 300; i++ {
			a := rand.Int() % 40
			b := rand.Int() % 40

			mulResult := functional.Execute[stand.MulTask, stand.MathInput, stand.MathOutput](
				ctx,
				stand.MulTask{},
				stand.MathInput{
					Left:  future.FromValue(a),
					Right: future.FromValue(b),
				},
			).Result

			result += a * b

			if i == 0 {
				previousFuture = mulResult
			} else {
				previousFuture = functional.Execute[stand.AddTask, stand.MathInput, stand.MathOutput](
					ctx,
					stand.AddTask{},
					stand.MathInput{
						Left:  previousFuture,
						Right: mulResult,
					},
				).Result
			}
		}
		return nil
	})
	if err != nil {
		panic(err)
	}

	fmt.Println("EXPECTED:", result)

	time.Sleep(time.Second * 4)

	for {
		time.Sleep(time.Second)
		fmt.Println("Expected:", result)

		tasks, err := recutil.List(
			ctx,
			sys.
				Tasks().
				Filter(record.R{"id": ops.In(taskCtx.GetSpawned()...)}).
				Cursor().
				Iter(),
		)
		if err != nil {
			fmt.Println("failed to load tasks:", err)
			continue
		}

		resources, err := sys.Resources().Load(ctx, lo.FlatMap(tasks, func(t task.Task, _ int) []resource.ID {
			return t.Outputs
		})...)
		if err != nil {
			fmt.Println("failed to load resources:", err)
			continue
		}

		resourceIndex := lo.SliceToMap(resources, func(res resource.Resource) (string, resource.Resource) {
			return res.ID, res
		})

		for _, t := range tasks {
			fmt.Println(t.ID)
			for index, out := range t.Outputs {
				val := resourceIndex[out]
				var repr string
				switch val.Status {
				case resource.Allocated:
					repr = "(allocated)"
				case resource.Ready:
					repr = string(val.Data)
				case resource.DoesNotExist:
					repr = "(none)"
				}
				fmt.Printf("\t%d => %s\n", index+1, repr)
			}
		}

	}
}

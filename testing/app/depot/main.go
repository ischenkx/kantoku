package main

import (
	"context"
	"kantoku"
	"kantoku/core/task"
	"kantoku/impl/common/pool/func"
	"kantoku/testing/app/base"
	"log"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	kan, err := base.Generate(ctx)
	if err != nil {
		log.Fatal(err)
		return
	}

	inputs := base.Inputs()

	scheduler := task.NewScheduler[kantoku.Task](inputs, kan.Events())
	err = kan.Depot().Process(
		ctx,
		funcpool.NewWriter[string](
			func(ctx context.Context, taskID string) error {
				task, err := kan.Tasks().Get(ctx, taskID)
				if err != nil {
					return err
				}
				return scheduler.Schedule(ctx, task)
			},
		),
	)
	if err != nil {
		log.Println("failed to run the depot processor:", err)
		return
	}
}

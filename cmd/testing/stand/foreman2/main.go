package main

import (
	"context"
	"fmt"
	"github.com/ischenkx/kantoku/cmd/testing/stand/common"
	"github.com/ischenkx/kantoku/pkg/processors/scheduler/dependencies/simple/manager"
	"github.com/ischenkx/kantoku/pkg/system"
	"github.com/ischenkx/kantoku/pkg/system/kernel/task"
	"time"
)

var Interval = time.Millisecond * 1000

func main() {
	common.InitLogger()
	ctx := context.Background()
	sys := common.NewSystem(ctx, "foreman-0")

	resources, err := sys.Resources().Alloc(ctx, 2)
	if err != nil {
		fmt.Println("failed to allocate resources:", err)
		return
	}

	in := resources[0]
	out := resources[1]

	fmt.Println("Input:", in)
	fmt.Println("Output:", out)

	t, err := sys.Spawn(ctx,
		system.WithInputs(in),
		system.WithOutputs(out),
		system.WithProperties(task.Properties{
			Data: map[string]any{
				"dependencies": []manager.DependencySpec{
					{
						Name: "resource",
						Data: in,
					},
				},
			},
		}),
	)

	if err != nil {
		fmt.Println("failed to spawn:", err)
		return
	}

	fmt.Println("task id:", t.ID)
}

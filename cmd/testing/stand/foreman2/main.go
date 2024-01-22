package main

import (
	"context"
	"fmt"
	"github.com/ischenkx/kantoku/cmd/testing/stand/common"
	"github.com/ischenkx/kantoku/pkg/common/data/record"
	"github.com/ischenkx/kantoku/pkg/core/resource"
	"github.com/ischenkx/kantoku/pkg/core/services/scheduler/dependencies/simple/manager"
	"github.com/ischenkx/kantoku/pkg/core/task"
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
		task.Task{
			Inputs:  []resource.ID{in},
			Outputs: []resource.ID{out},
			Info: record.R{
				"dependencies": []manager.DependencySpec{
					{
						Name: "resource",
						Data: in,
					},
				},
			},
		},
	)

	if err != nil {
		fmt.Println("failed to spawn:", err)
		return
	}

	fmt.Println("task id:", t.ID)
}

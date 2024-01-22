package main

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"github.com/ischenkx/kantoku/cmd/testing/stand/common"
	"github.com/ischenkx/kantoku/pkg/common/data/record"
	"github.com/ischenkx/kantoku/pkg/common/data/record/ops"
	"github.com/ischenkx/kantoku/pkg/core/resource"
	"github.com/ischenkx/kantoku/pkg/core/services/scheduler/dependencies/simple/manager"
	"github.com/ischenkx/kantoku/pkg/core/task"
	"github.com/samber/lo"
	"math/rand"
	"strconv"
)

const InitialInputsSize = 300
const TotalTasks = 180
const MinInputs = 1
const MaxInputs = 10
const MinOutputs = 1
const MaxOutputs = 10

func main() {
	common.InitLogger()
	ctx := context.Background()
	sys := common.NewSystem(ctx, "foreman-0")

	var inputs []string

	if initialInputs, err := sys.Resources().Alloc(ctx, InitialInputsSize); err != nil {
		fmt.Println("failed to allocate inputs:", err)
		return
	} else {
		_inputs := lo.Map(initialInputs, func(id string, _ int) resource.Resource {
			return resource.Resource{
				Data: []byte(strconv.Itoa(rand.Intn(1024))),
				ID:   id,
			}
		})

		if err := sys.Resources().Init(ctx, _inputs); err != nil {
			fmt.Println("failed to initialize inputs:", err)
			return
		}

		inputs = initialInputs
	}

	spawned := 0
	failed := 0

	var taskList []string

	for i := 0; i < TotalTasks; i++ {
		inputsAmount := MinInputs + rand.Intn(MaxInputs-MinInputs+1)
		outputsAmount := MinOutputs + rand.Intn(MaxOutputs-MinOutputs+1)

		inputIds := lo.Samples(inputs, inputsAmount)

		outputs, err := sys.Resources().Alloc(ctx, outputsAmount)
		if err != nil {
			fmt.Println("failed to allocate output resources:", err)
			failed++
			return
		}

		dependencies := lo.Map(inputIds, func(id string, _ int) manager.DependencySpec {
			return manager.DependencySpec{
				Name: "resource",
				Data: id,
			}
		})

		t, err := sys.Spawn(ctx,
			task.Task{
				Inputs:  inputIds,
				Outputs: outputs,
				Info: record.R{
					"dependencies": dependencies,
				},
			},
		)
		if err != nil {
			fmt.Println("failed to spawn:", err)
			failed++
			return
		}

		taskList = append(taskList, t.ID)
		inputs = append(inputs, outputs...)

		spawned++
	}

	fmt.Println("failed:", failed)
	fmt.Println("spawned:", spawned)

	for _, id := range taskList {
		fmt.Println("-", id)
	}

	ctxId := uuid.New().String()
	err := sys.
		Tasks().
		Filter(record.R{
			"id": ops.In(taskList...),
		}).
		Update(ctx, record.R{"context_id": ctxId}, nil)
	if err != nil {
		fmt.Println("failed to set context id:", err)
		return
	}

	fmt.Println("--------------")

	fmt.Println(ctxId)
	fmt.Printf("http://localhost:5173/sandbox/%s\n", ctxId)

	//ticker := time.NewTicker(time.Second)
	//defer ticker.Stop()

	//for range ticker.C {
	//	iter := sys.
	//		Info().
	//		Filter(record.R{
	//			system.InfoTaskID: ops.In(taskList...),
	//		}).
	//		Cursor().
	//		Iter()
	//
	//	list, err := recutil.ListIter[record.R](ctx, iter)
	//	if err != nil {
	//		fmt.Println("failed to load tasks:", err)
	//		continue
	//	}
	//
	//	ok := 0
	//	failed := 0
	//	pending := 0
	//	for _, t := range list {
	//		switch t["status"] {
	//		case "ok":
	//			ok++
	//		case "failed":
	//			failed++
	//		default:
	//			pending++
	//		}
	//	}
	//
	//	fmt.Printf("CTX=%s OK=%d FAILED=%d PENDING=%d\n", ctxId, ok, failed, pending)
	//}
}

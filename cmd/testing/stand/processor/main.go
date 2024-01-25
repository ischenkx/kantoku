package main

import (
	"context"
	"fmt"
	"github.com/ischenkx/kantoku/cmd/testing/stand/common"
	"github.com/ischenkx/kantoku/pkg/common/data/record"
	"github.com/ischenkx/kantoku/pkg/core/resource"
	"github.com/ischenkx/kantoku/pkg/core/task"
	"github.com/ischenkx/kantoku/pkg/lib/exe"
	"github.com/ischenkx/kantoku/pkg/lib/tasks/functional"
	"log"
	"log/slog"
	"strconv"
)

const Consumers = 5

//func execute(ctx *exe.Context) error {
//	slog.Info("executing", slog.String("id", ctx.Task().ID))
//
//	if len(ctx.Task().Inputs) != 1 {
//		return errors.New("expected one input: url")
//	}
//
//	if len(ctx.Task().Outputs) != 1 {
//		return errors.New("expected one output: response body")
//	}
//
//	urlResourceArray, err := ctx.System().Resources().Load(ctx, ctx.Task().Inputs[0])
//	if err != nil {
//		return fmt.Errorf("failed to load left resource: %w", err)
//	}
//
//	urlResource := urlResourceArray[0]
//	if urlResource.Status != resource.Ready {
//		return fmt.Errorf("url resource is not ready (id=%s)", urlResourceArray[0])
//	}
//
//	url := string(urlResource.Data)
//
//	slog.Debug("sending an http request",
//		slog.String("url", url))
//
//	response, err := http.Get(url)
//	if err != nil {
//		return fmt.Errorf("http request failed: %w", err)
//	}
//
//	slog.Debug("received left response",
//		slog.String("url", url),
//		slog.String("status", response.Status))
//
//	data, err := io.ReadAll(response.Body)
//	if err != nil {
//		return fmt.Errorf("failed to read the response body: %w", err)
//	}
//
//	err = ctx.System().Resources().Init(ctx, []resource.Resource{
//		{
//			Data: data,
//			ID:   ctx.Task().Outputs[0],
//		},
//	})
//	if err != nil {
//		return fmt.Errorf("failed save the output: %w", err)
//	}
//
//	return nil
//}

func execute(ctx *exe.Context) error {
	slog.Info("executing", slog.String("id", ctx.Task().ID))

	inputs, err := ctx.System().Resources().Load(ctx, ctx.Task().Inputs...)
	if err != nil {
		return fmt.Errorf("failed to load resources: %w", err)
	}

	outputs, err := ctx.System().Resources().Load(ctx, ctx.Task().Outputs...)
	if err != nil {
		return fmt.Errorf("failed to load resources: %w", err)
	}

	var numberInputs []int64
	for _, input := range inputs {
		num, err := strconv.ParseInt(string(input.Data), 10, 64)
		if err != nil {
			return fmt.Errorf("failed to parse an input: %w", err)
		}

		numberInputs = append(numberInputs, num)
	}

	typ, ok := ctx.Task().Info["type"]
	if !ok {
		typ = "sum"
	}

	for _, input := range inputs {
		if input.Status != resource.Ready {
			return fmt.Errorf("resource is not ready (id=%s)", input.ID)
		}
	}

	result := int64(0)
	for _, num := range numberInputs {
		switch typ {
		case "sum":
			result = result + num
		case "mul":
			result = result * num
		case "modbyte":
			result = (result + num) % 256
		case "modbit":
			result = (result + num) % 2
		}
	}

	for idx, out := range outputs {
		out.Data = []byte(strconv.Itoa(int(result * int64(idx+1))))
		outputs[idx] = out
	}

	slog.Info("done",
		slog.Any("type", typ),
		slog.Any("result", result))

	err = ctx.System().Resources().Init(ctx, outputs)
	if err != nil {
		return fmt.Errorf("failed save the outputs: %w", err)
	}

	return nil
}

func main() {
	common.InitLogger()

	slog.Info("Starting...")

	sys := common.NewSystem(context.Background(), "den-test")

	x1 := resource.Resource{
		Data:   []byte("1"),
		ID:     "init-res-1",
		Status: resource.Ready,
	}
	x2 := resource.Resource{
		Data:   []byte("2"),
		ID:     "init-res-2",
		Status: resource.Ready,
	}
	err := sys.Resources().Init(context.Background(), []resource.Resource{x1, x2})
	if err != nil {
		panic(err)
	}

	exec := functional.NewExecutor[AddTask, MathInput, MathOutput](AddTask{})
	err = exec.Execute(context.Background(), sys, task.Task{
		Inputs:  []resource.ID{x1.ID, x2.ID},
		Outputs: []resource.ID{"out-res-1"},
		ID:      "123",
		Info:    record.R{},
	})

	if err != nil {
		panic(err)
	}

	resources, err := sys.Resources().Load(context.Background(), "out-res-1")
	if err != nil {
		panic(err)
	}
	log.Println(string(resources[0].Data))
}

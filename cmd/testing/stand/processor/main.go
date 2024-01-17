package main

import (
	"context"
	"fmt"
	"github.com/ischenkx/kantoku/cmd/testing/stand/common"
	"github.com/ischenkx/kantoku/pkg/common/data/codec"
	"github.com/ischenkx/kantoku/pkg/extensions/exe"
	"github.com/ischenkx/kantoku/pkg/processors/executor"
	"github.com/ischenkx/kantoku/pkg/system/kernel/resource"
	"log/slog"
	"strconv"
	"sync"
)

const Consumers = 100

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
//		return fmt.Errorf("failed to load a resource: %w", err)
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
//	slog.Debug("received a response",
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

	typ, ok := ctx.Task().Properties.Data["type"]
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

	wg := &sync.WaitGroup{}

	for i := 0; i < Consumers; i++ {
		wg.Add(1)
		go func(index int) {
			defer wg.Done()

			proc := executor.NewProcessor(
				common.NewSystem(context.Background(), fmt.Sprintf("exe-%d", index)),
				exe.New(execute),
				"processor-1",
				codec.JSON[executor.Result](),
			)

			err := proc.Process(context.Background())
			if err != nil {
				slog.Error("failed:",
					slog.String("err", err.Error()))
				return
			}
		}(i)
	}

	slog.Info("Waiting...")
	wg.Wait()

	slog.Info("Finished!")
}

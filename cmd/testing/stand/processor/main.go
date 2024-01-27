package main

import (
	"context"
	"fmt"
	"github.com/ischenkx/kantoku/pkg/common/data/codec"
	"github.com/ischenkx/kantoku/pkg/common/service"
	"github.com/ischenkx/kantoku/pkg/core/resource"
	"github.com/ischenkx/kantoku/pkg/core/services/executor"
	"github.com/ischenkx/kantoku/pkg/lib/connector/cli/builder"
	config2 "github.com/ischenkx/kantoku/pkg/lib/connector/cli/config"
	"github.com/ischenkx/kantoku/pkg/lib/discovery"
	"github.com/ischenkx/kantoku/pkg/lib/exe"
	"github.com/joho/godotenv"
	"log"
	"log/slog"
	"strconv"
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

	slog.Debug("done",
		slog.Any("type", typ),
		slog.Any("result", result))

	err = ctx.System().Resources().Init(ctx, outputs)
	if err != nil {
		return fmt.Errorf("failed save the outputs: %w", err)
	}

	return nil
}

func main() {
	//common.InitLogger()

	if err := godotenv.Load("local/host.env"); err != nil {
		fmt.Println("failed to load env:", err)
		return
	}

	slog.Info("Starting...")

	var deployer service.Deployer

	cfg, err := config2.FromFile("local/config.yaml")
	if err != nil {
		log.Fatal(err)
	}
	var b builder.Builder
	for i := 0; i < Consumers; i++ {
		sys, err := b.BuildSystem(context.Background(), cfg.System)
		if err != nil {
			log.Fatal(err)
		}
		srvc := &executor.Service{
			System:      sys,
			ResultCodec: codec.JSON[executor.Result](),
			Executor:    exe.New(execute),
			Core: service.NewCore(
				"executor",
				fmt.Sprintf("exe-%d", i),
				slog.Default()),
		}
		deployer.Add(srvc,
			discovery.WithStaticInfo[*executor.Service](
				map[string]any{
					"executor": "simple",
				},
				sys.Events(),
				codec.JSON[discovery.Request](),
				codec.JSON[discovery.Response](),
			),
		)
	}

	if err := deployer.Deploy(context.Background()); err != nil {
		slog.Error("failed to deploy",
			slog.String("error", err.Error()))
	}
}

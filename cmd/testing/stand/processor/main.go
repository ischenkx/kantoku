package main

import (
	"context"
	"fmt"
	"github.com/ischenkx/kantoku/pkg/common/data/codec"
	"github.com/ischenkx/kantoku/pkg/common/service"
	"github.com/ischenkx/kantoku/pkg/core/services/executor"
	"github.com/ischenkx/kantoku/pkg/lib/connector/cli/builder"
	config2 "github.com/ischenkx/kantoku/pkg/lib/connector/cli/config"
	"github.com/ischenkx/kantoku/pkg/lib/discovery"
	"github.com/ischenkx/kantoku/pkg/lib/exe"
	"github.com/ischenkx/kantoku/pkg/lib/tasks/functional"
	"github.com/joho/godotenv"
	"log"
	"log/slog"
)

const Consumers = 5

func execute(ctx *exe.Context) error {
	err := functional.NewExecutor[AddTask, MathInput, MathOutput](AddTask{}).Execute(ctx, ctx.System(), ctx.Task())
	if err != nil {
		fmt.Println("failed to execute:", err)
		return nil
	}

	fmt.Println("executed:", ctx.Task().ID)

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

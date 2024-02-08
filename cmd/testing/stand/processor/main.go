package main

import (
	"context"
	"github.com/ischenkx/kantoku/cmd/testing/stand/common"
	"github.com/ischenkx/kantoku/pkg/common/data/codec"
	"github.com/ischenkx/kantoku/pkg/common/service"
	"github.com/ischenkx/kantoku/pkg/core/services/executor"
	"github.com/ischenkx/kantoku/pkg/lib/discovery"
	"github.com/ischenkx/kantoku/pkg/lib/exe"
	"github.com/ischenkx/kantoku/pkg/lib/tasks/functional"
	"log/slog"
)

func main() {
	common.InitLogger()

	slog.Info("Starting...")

	var deployer service.Deployer

	sys := common.NewSystem(context.Background(), "")

	router := exe.NewRouter()
	addExe := functional.NewExecutor[common.AddTask, common.MathInput, common.MathOutput](common.AddTask{})
	router.AddExecutor(addExe, addExe.Type())
	sumExe := functional.NewExecutor[common.SumTask, common.SumInput, common.MathOutput](common.SumTask{})
	router.AddExecutor(sumExe, sumExe.Type())

	srvc := &executor.Service{
		System:      sys,
		ResultCodec: codec.JSON[executor.Result](),
		Executor:    &router,
		Core: service.NewCore(
			"math-executor",
			"exe-math-0",
			slog.Default()),
	}
	deployer.Add(srvc,
		discovery.WithStaticInfo[*executor.Service](
			map[string]any{
				"executor": "functional",
			},
			sys.Events(),
			codec.JSON[discovery.Request](),
			codec.JSON[discovery.Response](),
		),
	)

	if err := deployer.Deploy(context.Background()); err != nil {
		slog.Error("failed to deploy",
			slog.String("error", err.Error()))
	}
}

package builder

import (
	"context"
	"github.com/ischenkx/kantoku/cmd/testing/stand/common"
	"github.com/ischenkx/kantoku/pkg/common/data/codec"
	"github.com/ischenkx/kantoku/pkg/common/service"
	"github.com/ischenkx/kantoku/pkg/core/services/executor"
	"github.com/ischenkx/kantoku/pkg/core/system"
	"github.com/ischenkx/kantoku/pkg/lib/connector/cli/config"
	"github.com/ischenkx/kantoku/pkg/lib/connector/cli/errx"
	"github.com/ischenkx/kantoku/pkg/lib/exe"
	"github.com/ischenkx/kantoku/pkg/lib/tasks/functional"
)

func (builder *Builder) BuildProcessor(ctx context.Context, sys system.AbstractSystem, cfg config.DynamicConfig) (service.Service, []service.Middleware, error) {
	router := exe.NewRouter()
	addExe := functional.NewExecutor[common.AddTask, common.MathInput, common.MathOutput](common.AddTask{})
	router.AddExecutor(addExe, addExe.Type())
	mulExe := functional.NewExecutor[common.MulTask, common.MathInput, common.MathOutput](common.MulTask{})
	router.AddExecutor(mulExe, mulExe.Type())
	divExe := functional.NewExecutor[common.DivTask, common.MathInput, common.MathOutput](common.DivTask{})
	router.AddExecutor(divExe, divExe.Type())
	sumExe := functional.NewExecutor[common.SumTask, common.SumInput, common.MathOutput](common.SumTask{})
	router.AddExecutor(sumExe, sumExe.Type())

	core, err := builder.BuildServiceCore(ctx, "status", cfg)
	if err != nil {
		return nil, nil, errx.FailedToBuild("core", err)
	}

	srvc := &executor.Service{
		System:      sys,
		ResultCodec: codec.JSON[executor.Result](),
		Executor:    &router,
		Core:        core,
	}

	mws, err := builder.BuildMiddlewares(ctx, srvc, sys, cfg)
	if err != nil {
		return nil, nil, errx.FailedToBuild("middlewares", err)
	}

	return srvc, mws, nil
}

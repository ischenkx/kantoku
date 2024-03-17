package builder

import (
	"context"
	"github.com/ischenkx/kantoku/cmd/testing/stand"
	"github.com/ischenkx/kantoku/pkg/common/data/codec"
	"github.com/ischenkx/kantoku/pkg/common/service"
	"github.com/ischenkx/kantoku/pkg/core/services/executor"
	"github.com/ischenkx/kantoku/pkg/core/system"
	"github.com/ischenkx/kantoku/pkg/lib/connector/cli/config"
	"github.com/ischenkx/kantoku/pkg/lib/connector/cli/errx"
)

func (builder *Builder) BuildProcessor(ctx context.Context, sys system.AbstractSystem, cfg config.DynamicConfig) (service.Service, []service.Middleware, error) {
	if cfg.Kind() != "math" {
		return nil, nil, errx.UnsupportedKind(cfg.Kind())
	}

	core, err := builder.BuildServiceCore(ctx, "processor", cfg)
	if err != nil {
		return nil, nil, errx.FailedToBuild("core", err)
	}

	srvc := &executor.Service{
		System:      sys,
		ResultCodec: codec.JSON[executor.Result](),
		Executor:    stand.MathExecutor(),
		Core:        core,
	}

	mws, err := builder.BuildMiddlewares(ctx, srvc, sys, cfg)
	if err != nil {
		return nil, nil, errx.FailedToBuild("middlewares", err)
	}

	return srvc, mws, nil
}

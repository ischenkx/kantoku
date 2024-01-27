package builder

import (
	"context"
	"github.com/ischenkx/kantoku/pkg/common/data/codec"
	"github.com/ischenkx/kantoku/pkg/common/service"
	"github.com/ischenkx/kantoku/pkg/core/services/executor"
	"github.com/ischenkx/kantoku/pkg/core/services/status"
	"github.com/ischenkx/kantoku/pkg/core/system"
	"github.com/ischenkx/kantoku/pkg/lib/connector/cli/config"
	"github.com/ischenkx/kantoku/pkg/lib/connector/cli/errx"
)

func (builder *Builder) BuildStatus(ctx context.Context, sys system.AbstractSystem, cfg config.DynamicConfig) (service.Service, []service.Middleware, error) {
	var statusConfig struct {
	}
	if err := cfg.Bind(&statusConfig); err != nil {
		return nil, nil, errx.FailedToBind(err)
	}

	core, err := builder.BuildServiceCore(ctx, "status", cfg)
	if err != nil {
		return nil, nil, errx.FailedToBuild("core", err)
	}

	ctx = withLogger(ctx, core.Logger())

	srvc := &status.Service{
		System:      sys,
		ResultCodec: codec.JSON[executor.Result](),
		Core:        core,
	}

	mws, err := builder.BuildMiddlewares(ctx, srvc, sys, cfg)
	if err != nil {
		return nil, nil, errx.FailedToBuild("middlewares", err)
	}

	return srvc, mws, nil
}

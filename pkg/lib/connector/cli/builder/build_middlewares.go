package builder

import (
	"context"
	"github.com/ischenkx/kantoku/pkg/common/data/codec"
	"github.com/ischenkx/kantoku/pkg/common/service"
	"github.com/ischenkx/kantoku/pkg/core/system"
	"github.com/ischenkx/kantoku/pkg/lib/connector/cli/config"
	"github.com/ischenkx/kantoku/pkg/lib/connector/cli/errx"
	"github.com/ischenkx/kantoku/pkg/lib/discovery"
	"golang.org/x/sync/errgroup"
	"log/slog"
)

func (builder *Builder) BuildMiddlewares(ctx context.Context, srvc service.Service, sys system.AbstractSystem, cfg config.DynamicConfig) (mws []service.Middleware, err error) {
	var data ServiceData[struct {
		Discovery struct {
			Enabled bool
			Info    map[string]any
		}
	}]
	if err := cfg.Bind(&data); err != nil {
		return nil, errx.FailedToBind(err)
	}

	if data.Info.Discovery.Enabled {
		mws = append(mws, discovery.WithStaticInfo[service.Service](
			data.Info.Discovery.Info,
			sys.Events(),
			codec.JSON[discovery.Request](),
			codec.JSON[discovery.Response](),
		))
	}

	mws = append(mws, loggingMiddleware{})

	return
}

type loggingMiddleware struct{}

func (l loggingMiddleware) BeforeRun(ctx context.Context, g *errgroup.Group, service service.Service) {
	service.Logger().Info("starting service",
		slog.String("name", service.Name()),
		slog.String("id", service.ID()))
}

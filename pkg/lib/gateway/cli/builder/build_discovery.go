package builder

import (
	"context"
	"fmt"
	capi "github.com/hashicorp/consul/api"
	"github.com/ischenkx/kantoku/pkg/common/data/codec"
	"github.com/ischenkx/kantoku/pkg/common/service"
	"github.com/ischenkx/kantoku/pkg/core/system"
	"github.com/ischenkx/kantoku/pkg/lib/discovery"
	"github.com/ischenkx/kantoku/pkg/lib/gateway/cli/config"
	"github.com/ischenkx/kantoku/pkg/lib/gateway/cli/errx"
	"github.com/ischenkx/kantoku/pkg/lib/impl/discovery/consul"
	"time"
)

func (builder *Builder) BuildDiscovery(ctx context.Context, sys system.AbstractSystem, cfg config.DynamicConfig) (service.Service, []service.Middleware, error) {
	var statusConfig struct {
		Hub             config.DynamicConfig
		PollingInterval time.Duration `mapstructure:"polling_interval"`
	}
	if err := cfg.Bind(&statusConfig); err != nil {
		return nil, nil, errx.FailedToBind(err)
	}

	core, err := builder.BuildServiceCore(ctx, "discovery", cfg)
	if err != nil {
		return nil, nil, errx.FailedToBuild("core", err)
	}

	ctx = withLogger(ctx, core.Logger())

	hub, err := builder.buildDiscoveryHub(ctx, statusConfig.Hub)
	if err != nil {
		return nil, nil, errx.FailedToBuild("hub", err)
	}

	srvc := &discovery.Poller{
		Hub:           hub,
		Events:        sys.Events(),
		RequestCodec:  codec.JSON[discovery.Request](),
		ResponseCodec: codec.JSON[discovery.Response](),
		Interval:      statusConfig.PollingInterval,
		Core:          core,
	}

	mws, err := builder.BuildMiddlewares(ctx, srvc, sys, cfg)
	if err != nil {
		return nil, nil, errx.FailedToBuild("middlewares", err)
	}

	return srvc, mws, nil
}

func (builder *Builder) buildDiscoveryHub(ctx context.Context, cfg config.DynamicConfig) (discovery.Hub, error) {
	switch cfg.Kind() {
	case "consul":
		return builder.buildConsulHub(ctx, cfg)
	default:
		return nil, errx.UnsupportedKind(cfg.Kind())
	}
}

func (builder *Builder) buildConsulHub(ctx context.Context, cfg config.DynamicConfig) (*consul.Hub, error) {
	var consulConfig struct {
		Addr string
	}
	if err := cfg.Bind(&consulConfig); err != nil {
		return nil, errx.FailedToBind(err)
	}

	client, err := capi.NewClient(&capi.Config{
		Address: consulConfig.Addr,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create a consul client: %w", err)
	}

	return &consul.Hub{Consul: client}, nil
}

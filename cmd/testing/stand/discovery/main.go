package main

import (
	"context"
	"github.com/hashicorp/consul/api"
	"github.com/ischenkx/kantoku/cmd/testing/stand/common"
	"github.com/ischenkx/kantoku/pkg/common/data/codec"
	"github.com/ischenkx/kantoku/pkg/common/data/uid"
	"github.com/ischenkx/kantoku/pkg/common/service"
	"github.com/ischenkx/kantoku/pkg/lib/discovery"
	"github.com/ischenkx/kantoku/pkg/lib/impl/discovery/consul"
	"log/slog"
	"time"
)

func main() {
	common.InitLogger()
	ctx := context.Background()

	system := common.NewSystem(context.Background(), "scheduler-0")

	consulClient, err := api.NewClient(api.DefaultConfig())
	if err != nil {
		slog.Error("failed to connect to consul",
			slog.String("error", err.Error()))

		return
	}

	hub := &consul.Hub{
		Consul: consulClient,
	}

	srvc := &discovery.Poller{
		Hub:           hub,
		Events:        system.Events(),
		RequestCodec:  codec.JSON[discovery.Request](),
		ResponseCodec: codec.JSON[discovery.Response](),
		Interval:      time.Second * 3,
		Core: service.NewCore(
			"discovery-poller",
			uid.Generate(),
			slog.Default(),
		),
	}

	var deployer service.Deployer

	deployer.Add(
		srvc,
		discovery.WithStaticInfo[*discovery.Poller](
			map[string]any{},
			system.Events(),
			codec.JSON[discovery.Request](),
			codec.JSON[discovery.Response](),
		))

	if err := deployer.Deploy(ctx); err != nil {
		slog.Error("failed to deploy",
			slog.String("error", err.Error()))
	}
}

type _hub struct {
}

func (h _hub) Register(ctx context.Context, info discovery.ServiceInfo) error {
	slog.Info("registering",
		slog.String("service", info.Name),
		slog.String("id", info.ID))

	return nil
}

func (h _hub) Load(ctx context.Context) ([]discovery.ServiceInfo, error) {
	//TODO implement me
	panic("implement me")
}

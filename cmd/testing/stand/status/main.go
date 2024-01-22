package main

import (
	"context"
	"github.com/ischenkx/kantoku/cmd/testing/stand/common"
	codec "github.com/ischenkx/kantoku/pkg/common/data/codec"
	"github.com/ischenkx/kantoku/pkg/common/data/uid"
	"github.com/ischenkx/kantoku/pkg/common/service"
	"github.com/ischenkx/kantoku/pkg/core/services/executor"
	"github.com/ischenkx/kantoku/pkg/core/services/status"
	"github.com/ischenkx/kantoku/pkg/lib/discovery"

	"log/slog"
)

func main() {
	common.InitLogger()

	system := common.NewSystem(context.Background(), "status-0")

	srvc := &status.Service{
		System:      system,
		ResultCodec: codec.JSON[executor.Result](),
		Core: service.NewCore(
			"status",
			uid.Generate(),
			slog.Default(),
		),
	}

	var deployer service.Deployer

	deployer.Add(
		srvc,
		discovery.WithStaticInfo[*status.Service](
			map[string]any{},
			system.Events(),
			codec.JSON[discovery.Request](),
			codec.JSON[discovery.Response](),
		))

	if err := deployer.Deploy(context.Background()); err != nil {
		slog.Error("failed to deploy",
			slog.String("error", err.Error()))
	}
}

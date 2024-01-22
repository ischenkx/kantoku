package main

import (
	"context"
	"github.com/ischenkx/kantoku/cmd/testing/stand/common"
	"github.com/ischenkx/kantoku/pkg/common/data/codec"
	"github.com/ischenkx/kantoku/pkg/common/data/uid"
	"github.com/ischenkx/kantoku/pkg/common/service"
	"github.com/ischenkx/kantoku/pkg/core/system"
	"github.com/ischenkx/kantoku/pkg/lib/connector/api/http/oas"
	"github.com/ischenkx/kantoku/pkg/lib/connector/api/http/server"
	"github.com/ischenkx/kantoku/pkg/lib/discovery"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"log/slog"
)

func main() {
	ctx := context.Background()

	sys := common.NewSystem(ctx, "http-server-0")

	err := service.NewDeployer().
		Add(
			serverService{
				sys:  sys,
				Core: service.NewCore("http_server", uid.Generate(), slog.Default()),
			},
			discovery.WithStaticInfo[serverService](
				map[string]any{},
				sys.Events(),
				codec.JSON[discovery.Request](),
				codec.JSON[discovery.Response](),
			),
		).
		Deploy(context.Background())
	if err != nil {
		slog.Error("failed to deploy",
			slog.String("error", err.Error()))
	}
}

type serverService struct {
	sys system.AbstractSystem
	service.Core
}

func (srvc serverService) Run(ctx context.Context) error {
	srv := server.New(srvc.sys)

	e := echo.New()

	// Middleware
	e.Use(middleware.Logger())

	oas.RegisterHandlers(e, oas.NewStrictHandler(srv, nil))

	if err := e.Start(":8080"); err != nil {
		return err
	}

	return nil
}

package platform

import (
	"context"
	"fmt"
	"github.com/ischenkx/kantoku/pkg/common/service"
	"github.com/ischenkx/kantoku/pkg/core/system"
	"github.com/ischenkx/kantoku/pkg/lib/gateway/api/http"
	"github.com/ischenkx/kantoku/pkg/lib/gateway/api/http/oas"
	"github.com/ischenkx/kantoku/pkg/lib/tasks/specification"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"golang.org/x/sync/errgroup"
	"log/slog"
)

type Deployment[S service.Service] struct {
	Service     S
	Middlewares []service.Middleware
}

type HttpApiService struct {
	sys            system.AbstractSystem
	specifications *specification.Manager
	port           int
	loggerEnabled  bool
	service.Core
}

func (srvc HttpApiService) Run(ctx context.Context) error {
	srv := http.NewServer(
		srvc.sys,
		srvc.specifications,
	)

	e := echo.New()

	if srvc.loggerEnabled {
		e.Use(middleware.Logger())
	}
	e.Use(middleware.CORS())

	oas.RegisterHandlers(e, oas.NewStrictHandler(srv, nil))

	if err := e.Start(fmt.Sprintf(":%d", srvc.port)); err != nil {
		return err
	}

	return nil
}

type loggingMiddleware struct{}

func (l loggingMiddleware) BeforeRun(ctx context.Context, g *errgroup.Group, service service.Service) {
	service.Logger().Info("starting service",
		slog.String("name", service.Name()),
		slog.String("id", service.ID()))
}

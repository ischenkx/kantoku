package builder

import (
	"context"
	"fmt"
	"github.com/ischenkx/kantoku/pkg/common/service"
	"github.com/ischenkx/kantoku/pkg/core/system"
	"github.com/ischenkx/kantoku/pkg/lib/gateway/api/http"
	"github.com/ischenkx/kantoku/pkg/lib/gateway/api/http/oas"
	"github.com/ischenkx/kantoku/pkg/lib/gateway/cli/config"
	"github.com/ischenkx/kantoku/pkg/lib/gateway/cli/errx"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"strconv"
)

func (builder *Builder) BuildHttpApi(ctx context.Context, sys system.AbstractSystem, cfg config.DynamicConfig) (service.Service, []service.Middleware, error) {
	var apiConfig struct {
		Port          string
		LoggerEnabled bool `mapstructure:"logger_enabled"`
	}
	if err := cfg.Bind(&apiConfig); err != nil {
		return nil, nil, errx.FailedToBind(err)
	}

	core, err := builder.BuildServiceCore(ctx, "http-api", cfg)
	if err != nil {
		return nil, nil, errx.FailedToBuild("core", err)
	}

	ctx = withLogger(ctx, core.Logger())

	port, err := strconv.Atoi(apiConfig.Port)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to parse port (value='%s'): %w", apiConfig.Port, err)
	}

	srvc := httpApiService{
		sys:           sys,
		port:          port,
		loggerEnabled: apiConfig.LoggerEnabled,
		Core:          core,
	}

	mws, err := builder.BuildMiddlewares(ctx, srvc, sys, cfg)
	if err != nil {
		return nil, nil, errx.FailedToBuild("middlewares", err)
	}

	return srvc, mws, nil
}

type httpApiService struct {
	sys           system.AbstractSystem
	port          int
	loggerEnabled bool
	service.Core
}

func (srvc httpApiService) Run(ctx context.Context) error {
	srv := http.New(srvc.sys)

	e := echo.New()

	if srvc.loggerEnabled {
		e.Use(middleware.Logger())
	}

	oas.RegisterHandlers(e, oas.NewStrictHandler(srv, nil))

	if err := e.Start(fmt.Sprintf(":%d", srvc.port)); err != nil {
		return err
	}

	return nil
}

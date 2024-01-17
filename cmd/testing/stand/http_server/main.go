package main

import (
	"context"
	"fmt"
	"github.com/ischenkx/kantoku/cmd/testing/stand/common"
	"github.com/ischenkx/kantoku/pkg/extensions/api/http/oas"
	"github.com/ischenkx/kantoku/pkg/extensions/api/http/server"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	ctx := context.Background()
	sys := common.NewSystem(ctx, "http-server-0")
	srv := server.New(sys)

	e := echo.New()

	// Middleware
	e.Use(middleware.Logger())

	oas.RegisterHandlers(e, oas.NewStrictHandler(srv, nil))

	if err := e.Start(":8080"); err != nil {
		fmt.Println("failed:", err)
	}
}

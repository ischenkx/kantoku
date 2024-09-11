package main

import (
	"context"
	"github.com/ischenkx/kantoku/cmd/stand/utils"
	"github.com/ischenkx/kantoku/pkg/common/service"
	"github.com/ischenkx/kantoku/pkg/lib/builder"
	"log"
	"os"
)

func main() {
	ctx := context.Background()
	cfg := utils.LoadConfig()
	logger := utils.GetLogger(os.Stdout, "http_api")

	sys, err := builder.BuildSystem(ctx, logger, cfg.Core.System)
	if err != nil {
		log.Fatal("failed to build system: ", err)
	}

	specifications, err := builder.BuildSpecifications(ctx, cfg.Core.Specifications)
	if err != nil {
		log.Fatal("failed to build specifications:", err)
	}

	deployment, err := builder.BuildHttpApiDeployment(ctx, sys, specifications, logger, cfg.Services.HttpApi)
	if err != nil {
		log.Fatal("failed to build http api:", err)
	}

	deployer := service.NewDeployer()
	deployer.Add(deployment.Service, deployment.Middlewares...)
	if err := deployer.Deploy(ctx); err != nil {
		log.Fatal("failed to deploy:", err)
	}
}

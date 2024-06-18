package main

import (
	"context"
	"github.com/ischenkx/kantoku/cmd/stand/utils"
	"github.com/ischenkx/kantoku/pkg/common/service"
	"github.com/ischenkx/kantoku/pkg/lib/platform"
	"log"
	"os"
)

func main() {
	ctx := context.Background()
	cfg := utils.LoadConfig()
	logger := utils.GetLogger(os.Stdout, "scheduler")

	sys, err := platform.BuildSystem(ctx, logger, cfg.Core.System)
	if err != nil {
		log.Fatal("failed to build system: ", err)
	}

	deployment, err := platform.BuildSchedulerDeployment(ctx, sys, logger, cfg.Services.Scheduler)
	if err != nil {
		log.Fatal("failed to build scheduler:", err)
	}

	deployer := service.NewDeployer()
	deployer.Add(deployment.Service, deployment.Middlewares...)
	if err := deployer.Deploy(ctx); err != nil {
		log.Fatal("failed to deploy:", err)
	}
}

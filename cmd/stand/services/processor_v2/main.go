package main

import (
	"context"
	"github.com/ischenkx/kantoku/cmd/stand/executor"
	"github.com/ischenkx/kantoku/cmd/stand/utils"
	"github.com/ischenkx/kantoku/pkg/common/service"
	"github.com/ischenkx/kantoku/pkg/lib/platform"
	"log"
	"os"
)

func main() {
	ctx := context.Background()
	cfg := utils.LoadConfig()
	logger := utils.GetLogger(os.Stdout, "processor_v2")

	sys, err := platform.BuildSystem(ctx, logger, cfg.Core.System)
	if err != nil {
		log.Fatal("failed to build system: ", err)
	}

	exec := executor.New()

	deployment, err := platform.BuildProcessorDeployment(ctx, sys, exec, logger, cfg.Services.Processor)
	if err != nil {
		log.Fatal("failed to build processor:", err)
	}

	deployer := service.NewDeployer()
	deployer.Add(deployment.Service, deployment.Middlewares...)
	if err := deployer.Deploy(ctx); err != nil {
		log.Fatal("failed to deploy:", err)
	}
}

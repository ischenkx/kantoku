package main

import (
	"context"
	"github.com/ischenkx/kantoku/cmd/stand/utils"
	"github.com/ischenkx/kantoku/local/deprecated/math_executor"
	"github.com/ischenkx/kantoku/pkg/common/service"
	"github.com/ischenkx/kantoku/pkg/lib/platform"
	"log"
	"os"
)

func main() {
	ctx := context.Background()
	cfg := utils.LoadConfig()
	logger := utils.GetLogger(os.Stdout, "processor")

	sys, err := platform.BuildSystem(ctx, logger, cfg.Core.System)
	if err != nil {
		log.Fatal("failed to build system: ", err)
	}

	deployment, err := platform.BuildProcessorDeployment(ctx, sys, math_executor.MathExecutor(), logger, cfg.Services.Processor)
	if err != nil {
		log.Fatal("failed to build processor:", err)
	}

	deployer := service.NewDeployer()
	deployer.Add(deployment.Service, deployment.Middlewares...)
	if err := deployer.Deploy(ctx); err != nil {
		log.Fatal("failed to deploy:", err)
	}
}

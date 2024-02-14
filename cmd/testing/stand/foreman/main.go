package main

import (
	"context"
	"github.com/ischenkx/kantoku/pkg/core/task"
	"github.com/ischenkx/kantoku/pkg/lib/connector/cli/builder"
	"github.com/ischenkx/kantoku/pkg/lib/connector/cli/config"
	taskutil "github.com/ischenkx/kantoku/pkg/lib/tasks/util"
	"log"
)

func main() {
	ctx := context.Background()

	cfg, err := config.FromFile("local/config.yaml")
	if err != nil {
		log.Fatal(err)
	}
	var b builder.Builder

	sys, err := b.BuildSystem(ctx, cfg.System)
	if err != nil {
		log.Fatal(err)
	}

	sys.Spawn(ctx, task.New(
		task.WithInputs("1", "2"),
		task.WithOutputs("1", "2"),
		taskutil.DependOnInputs(),
	))
}

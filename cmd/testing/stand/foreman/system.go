package main

import (
	"context"
	"github.com/ischenkx/kantoku/pkg/core/system"
	"github.com/ischenkx/kantoku/pkg/lib/connector/cli/builder"
	config2 "github.com/ischenkx/kantoku/pkg/lib/connector/cli/config"
	"github.com/joho/godotenv"
)

func NewSystem(ctx context.Context) system.AbstractSystem {
	if err := godotenv.Load("local/host.env"); err != nil {
		panic(err)
	}

	cfg, err := config2.FromFile("local/config.yaml")

	myBuilder := builder.Builder{}
	sys, err := myBuilder.BuildSystem(ctx, cfg.System)
	if err != nil {
		panic(err)
	}
	return sys
}

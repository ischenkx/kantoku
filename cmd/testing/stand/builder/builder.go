package builder

import (
	"context"
	"github.com/ischenkx/kantoku/pkg/core/system"
	"github.com/ischenkx/kantoku/pkg/lib/gateway/cli/builder"
	config2 "github.com/ischenkx/kantoku/pkg/lib/gateway/cli/config"
	"github.com/joho/godotenv"
)

func NewSystem(ctx context.Context) system.AbstractSystem {
	bldr, cfg := NewBuilder(ctx)
	sys, err := bldr.BuildSystem(ctx, cfg.System)
	if err != nil {
		panic(err)
	}
	return sys
}

func NewBuilder(ctx context.Context) (builder.Builder, config2.Config) {
	if err := godotenv.Load("local/host.env"); err != nil {
		panic(err)
	}

	cfg, err := config2.FromFile("local/config.yaml")
	if err != nil {
		panic(err)
	}

	myBuilder := builder.Builder{}
	return myBuilder, cfg
}

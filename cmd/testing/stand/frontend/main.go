package main

import (
	"context"
	"github.com/ischenkx/kantoku/cmd/testing/stand/builder"
)

func main() {
	ctx := context.Background()
	bldr, cfg := builder.NewBuilder(ctx)
	api, _, err := bldr.BuildHttpApi(ctx, builder.NewSystem(ctx), cfg.Services.HttpServer)
	if err != nil {
		panic(err)
	}

	err = api.Run(ctx)
	if err != nil {
		panic(err)
	}
}

package main

import (
	"context"
	"kantoku/app/base"
	"kantoku/common/util"
	"kantoku/impl/common/deps/postgredeps"
	"log"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	kan, err := base.Generate(ctx)
	if err != nil {
		log.Fatal(err)
		return
	}

	deps := kan.Depot().Deps().(*postgredeps.Deps)
	deps.Run(ctx)

	util.Wait(ctx)
}

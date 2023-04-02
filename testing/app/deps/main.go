package main

import (
	"context"
	"kantoku/common/util"
	"kantoku/testing/app/base"
	"log"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	b, err := base.New(ctx)
	if err != nil {
		log.Fatal(err)
	}

	b.Deps.Run(ctx)
	util.BlockOn(ctx)
}

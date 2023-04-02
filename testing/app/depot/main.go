package main

import (
	"context"
	"kantoku/impl/common/codec/strcodec"
	redipool "kantoku/impl/common/pool/redis"
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

	outputs := redipool.New[string](b.Redis, strcodec.Codec{}, "tasks")

	if err := b.Depot.Process(ctx, outputs); err != nil {
		log.Fatal(err)
	}
}

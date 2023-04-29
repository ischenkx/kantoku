package main

import (
	"context"
	delay2 "kantoku/framework/depot/delay"
	"kantoku/impl/common/codec/jsoncodec"
	"kantoku/impl/common/codec/strcodec"
	"kantoku/impl/common/data/cron/simple"
	redipool "kantoku/impl/common/data/pool/redis"
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

	cronInputs := redipool.New[simple.Event](b.Redis, jsoncodec.New[simple.Event](), "cron_inputs")
	cronOutputs := redipool.New[string](b.Redis, strcodec.Codec{}, "cron_outputs")
	redisCron := simple.NewClient(cronInputs, cronOutputs)
	delayManager := delay2.NewManager(redisCron, b.Deps)

	if err := delay2.NewUpdater(delayManager).Run(ctx); err != nil {
		log.Fatal(err)
	}
}

package main

import (
	"context"
	"github.com/go-co-op/gocron"
	"kantoku/impl/common/codec/jsoncodec"
	"kantoku/impl/common/codec/strcodec"
	"kantoku/impl/common/data/cron/simple"
	redipool "kantoku/impl/common/data/pool/redis"
	"kantoku/testing/app/base"
	"log"
	"time"
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
	scheduler := gocron.NewScheduler(time.Now().Location())
	server := simple.NewServer(scheduler, cronInputs, cronOutputs)
	if err := server.Run(ctx); err != nil {
		log.Fatal(err)
	}
}

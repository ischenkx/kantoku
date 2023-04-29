package main

import (
	"context"
	"kantoku/framework/status"
	"kantoku/impl/common/codec/jsoncodec"
	redikv "kantoku/impl/common/data/kv/redis"
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

	statusDB := redikv.New[status.Status](b.Redis, jsoncodec.New[status.Status](), "statuses")

	status.NewUpdater(b.Kantoku.Events(), statusDB).Run(ctx)
}

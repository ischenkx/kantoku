package main

import (
	"context"
	subtask2 "kantoku/framework/depot/taskdep"
	"kantoku/impl/common/codec/strcodec"
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

	manager := subtask2.NewManager(b.Deps, redikv.New[string](b.Redis, strcodec.Codec{}, "subtasks"))
	if err := subtask2.NewUpdater(b.Kantoku.Events(), manager).Run(ctx); err != nil {
		log.Fatal(err)
	}
}

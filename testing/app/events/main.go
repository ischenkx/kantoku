package main

import (
	"context"
	"kantoku/core/task"
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

	ch, err := b.Kantoku.Events().Listen(ctx, task.EventTopic)
	if err != nil {
		log.Fatal(err)
	}

	for event := range ch {
		log.Printf("%s: %s", event.Name, event.Data)
	}
}
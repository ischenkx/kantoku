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

	kan, err := base.Generate(ctx)
	if err != nil {
		log.Fatal(err)
		return
	}

	events, err := kan.Events().Listen(ctx, task.EventTopic)
	if err != nil {
		log.Fatal(err)
		return
	}
	log.Println("Events:")
	for ev := range events {
		log.Printf("Name: '%s', Topic: '%s', Data: '%s'\n", ev.Name, ev.Topic, string(ev.Data))
	}
}

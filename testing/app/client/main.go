package main

import (
	"context"
	"kantoku"
	"kantoku/core/l1"
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
	log.Println("Starting...")

	id, err := kan.New(ctx, kantoku.Task{
		Spec: l1.Task{
			Type:     "reverse",
			Argument: []byte("Hello World!"),
		},
		Dependencies: nil,
	})
	if err != nil {
		log.Fatal(err)
		return
	}

	log.Println("task id:", id)

	log.Println("Done, enjoy!")
}

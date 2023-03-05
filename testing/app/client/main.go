package main

import (
	"context"
	"github.com/google/uuid"
	"kantoku"
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

	dep := uuid.New().String()

	id, err := kan.New(ctx, kantoku.Task{
		Type_:        "reverse",
		Argument_:    []byte("Hello World!"),
		Dependencies: []string{dep},
	})
	if err != nil {
		log.Fatal(err)
		return
	}

	log.Println("task id:", id)
	log.Println("dep:", dep)
	log.Println("Done, enjoy!")
}

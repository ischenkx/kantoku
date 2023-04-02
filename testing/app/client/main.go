package main

import (
	"context"
	"kantoku"
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

	_, err = b.Kantoku.Spawn(context.Background(),
		kantoku.Task("http", "https://google.com"),
	)

	if err != nil {
		log.Fatal(err)
	}

	log.Println("Done, enjoy!")
}

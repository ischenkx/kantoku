package main

import (
	"context"
	"kantoku"
	"kantoku/framework/depot/delay"
	"kantoku/framework/depot/taskdep"
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

	result1, err := b.Kantoku.Spawn(
		context.Background(),
		kantoku.Task("fibonacci", []byte("10")).
			With(delay.Delay(time.Second*5)),
	)

	if err != nil {
		log.Fatal(err)
	}

	result2, err := b.Kantoku.Spawn(
		context.Background(),
		kantoku.Task("fibonacci", []byte("10")).
			With(delay.Delay(time.Second*5)).
			With(taskdep.Dep(result1.Task)),
	)

	log.Println("Task1:", result1.Task)
	log.Println("Task2:", result2.Task)
	log.Println("Done, enjoy!")
}

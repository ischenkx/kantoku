package main

import (
	"context"
	"github.com/redis/go-redis/v9"
	"github.com/samber/lo"
	"kantoku"
	"kantoku/app/base"
	"kantoku/core/task"
	"kantoku/framework/executors/simple"
	"kantoku/impl/common/codec/jsoncodec"
	redikv "kantoku/impl/common/data/kv/redis"
	"log"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	redisClient := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})

	kan, err := base.Generate(ctx)
	if err != nil {
		log.Fatal(err)
		return
	}

	executor := simple.Executor[kantoku.Task]{
		"reverse": func(data []byte) ([]byte, error) {
			return lo.Reverse(data), nil
		},
	}

	inputs := base.Inputs()
	outputs := base.Outputs(redikv.New[task.Result](redisClient, jsoncodec.New[task.Result](), "outputs"))

	if err := task.NewPipeline[kantoku.Task](inputs, outputs, executor, kan.Events()).Run(context.Background()); err != nil {
		log.Println("failed to run the task runner:", err)
	}
}

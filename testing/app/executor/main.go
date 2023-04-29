package main

import (
	"context"
	"kantoku"
	"kantoku/common/data/kv"
	"kantoku/common/data/pool"
	"kantoku/core/task"
	"kantoku/framework/executors/simple"
	"kantoku/impl/common/codec/jsoncodec"
	"kantoku/impl/common/codec/strcodec"
	redikv "kantoku/impl/common/data/kv/redis"
	"kantoku/impl/common/data/pool/proxypool"
	"kantoku/impl/common/data/pool/redis"
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

	ids := redipool.New[string](b.Redis, strcodec.Codec{}, "_tasks")
	inputs := newInputs(ids, b.Kantoku)
	outputs := redikv.New[task.Result](b.Redis, jsoncodec.New[task.Result](), "outputs")

	p := task.NewPipeline[*kantoku.View](inputs, Outputs{storage: outputs}, executor(), b.Kantoku.Events())
	if err := p.Run(ctx); err != nil {
		log.Fatal(err)
	}
}

func executor() task.Executor[*kantoku.View] {
	return simple.Executor{
		"print": func(ctx context.Context, task *kantoku.View) ([]byte, error) {
			data, err := task.Data(ctx)
			if err != nil {
				return nil, err
			}
			log.Println(string(data))

			return nil, nil
		},
		"id": func(ctx context.Context, task *kantoku.View) ([]byte, error) {
			return task.Data(ctx)
		},
	}
}

func newInputs(ids pool.Reader[string], kan *kantoku.Kantoku) pool.Reader[*kantoku.View] {
	return proxypool.NewReader[string, *kantoku.View](ids, func(ctx context.Context, id string) (*kantoku.View, bool) {
		return kan.Task(id), true
	})
}

type Outputs struct {
	storage kv.Database[string, task.Result]
}

func (o Outputs) Write(ctx context.Context, item task.Result) error {
	return o.storage.Set(ctx, item.TaskID, item)
}

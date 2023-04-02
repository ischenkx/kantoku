package main

import (
	"context"
	"errors"
	"io"
	"kantoku"
	"kantoku/common/data/kv"
	"kantoku/common/data/pool"
	"kantoku/core/task"
	"kantoku/framework/executors/simple"
	"kantoku/impl/common/codec/jsoncodec"
	"kantoku/impl/common/codec/strcodec"
	redikv "kantoku/impl/common/data/kv/redis"
	"kantoku/impl/common/pool/proxypool"
	redipool "kantoku/impl/common/pool/redis"
	"kantoku/testing/app/base"
	"log"
	"net/http"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	b, err := base.New(ctx)
	if err != nil {
		log.Fatal(err)
	}

	ids := redipool.New[string](b.Redis, strcodec.Codec{}, "tasks")
	inputs := newInputs(ids, b.Kantoku)
	outputs := redikv.New[task.Result](b.Redis, jsoncodec.New[task.Result](), "outputs")

	p := task.NewPipeline[kantoku.StoredTask](inputs, Outputs{storage: outputs}, executor(), b.Kantoku.Events())
	if err := p.Run(ctx); err != nil {
		log.Fatal(err)
	}
}

func executor() task.Executor[kantoku.StoredTask] {
	return simple.Executor{
		"http": func(argument any) ([]byte, error) {
			url, ok := argument.(string)
			if !ok {
				return nil, errors.New("failed to get the url")
			}

			res, err := http.Get(url)
			if err != nil {
				return nil, err
			}

			return io.ReadAll(res.Body)
		},
	}
}

func newInputs(ids pool.Reader[string], kan *kantoku.Kantoku) pool.Reader[kantoku.StoredTask] {
	return proxypool.NewReader[string, kantoku.StoredTask](ids, func(ctx context.Context, id string) (kantoku.StoredTask, bool) {
		stored, err := kan.Task(id).AsStored(ctx)
		if err != nil {
			log.Println("failed to get the stored task:", err)
			return kantoku.StoredTask{}, false
		}
		return stored, true
	})
}

type Outputs struct {
	storage kv.Database[task.Result]
}

func (o Outputs) Write(ctx context.Context, item task.Result) error {
	_, err := o.storage.Set(ctx, item.TaskID, item)
	return err
}

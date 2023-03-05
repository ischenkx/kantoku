package base

import (
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
	"kantoku"
	"kantoku/common/data/kv"
	"kantoku/common/pool"
	"kantoku/core/event"
	"kantoku/core/task"
	"kantoku/framework/depot"
	"kantoku/impl/common/codec/bincodec"
	"kantoku/impl/common/codec/jsoncodec"
	"kantoku/impl/common/codec/strcodec"
	redikv "kantoku/impl/common/data/kv/redis"
	"kantoku/impl/common/deps/postgredeps"
	funcpool "kantoku/impl/common/pool/func"
	redipool "kantoku/impl/common/pool/redis"
	"kantoku/impl/core/event/redis"
	"kantoku/impl/framework/cell/redis"
)

type IdentifiableResult task.Result

func (i IdentifiableResult) ID() string {
	return i.TaskID
}

func Inputs() pool.Pool[kantoku.Task] {
	redisClient := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})
	return redipool.New[kantoku.Task](redisClient, jsoncodec.New[kantoku.Task](), "inputs")
}

func Outputs(outputs kv.Writer[task.Result]) pool.Writer[task.Result] {
	return funcpool.NewWriter[task.Result](
		func(ctx context.Context, item task.Result) error {
			_, err := outputs.Set(ctx, item.TaskID, item)
			return err
		},
	)
}

func Generate(ctx context.Context) (*kantoku.Kantoku, error) {
	redisClient := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})

	postgres, err := pgxpool.New(ctx, "postgresql://postgres:root@localhost:5432/postgres")
	if err != nil {
		return nil, err
	}

	bus := redivent.NewBus(redisClient, jsoncodec.New[event.Event]())
	depsQueue := redipool.New[string](redisClient, strcodec.Codec{}, "deps_queue")
	deps := postgredeps.New(postgres, depsQueue)
	tasks := redikv.New[kantoku.Task](redisClient, jsoncodec.New[kantoku.Task](), "tasks")
	cells := redicell.New[[]byte](redisClient, bincodec.Codec{})

	dep := depot.New(
		deps,
		redikv.New[string](redisClient, jsoncodec.New[string](), "deps_mapper"),
	)

	return kantoku.New(
		kantoku.Config{
			Events: bus,
			Depot:  dep,
			Tasks:  tasks,
			Cells:  cells,
		},
	), nil
}

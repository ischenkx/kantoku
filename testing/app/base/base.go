package base

import (
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
	"kantoku"
	"kantoku/common/pool"
	"kantoku/core/l0/event"
	"kantoku/core/l1"
	"kantoku/framework/cell"
	"kantoku/framework/depot"
	"kantoku/impl/common/codec/bincodec"
	"kantoku/impl/common/codec/jsoncodec"
	"kantoku/impl/common/codec/strcodec"
	redikv "kantoku/impl/common/db/kv/redis"
	"kantoku/impl/common/deps/postgredeps"
	funcpool "kantoku/impl/common/pool/func"
	redipool "kantoku/impl/common/pool/redis"
	redivent "kantoku/impl/core/l0/event/redis"
	redicell "kantoku/impl/framework/cell/redis"
)

type IdentifiableResult l1.Result

func (i IdentifiableResult) ID() string {
	return i.TaskID
}

func L1Inputs() pool.Pool[l1.Task] {
	redisClient := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})
	return redipool.New[l1.Task](redisClient, jsoncodec.New[l1.Task](), "inputs")
}

func L1Outputs(cells cell.Storage[[]byte]) pool.Writer[l1.Result] {
	return funcpool.NewWriter[l1.Result](
		func(ctx context.Context, item l1.Result) error {
			return cells.Set(ctx, item.TaskID, item.Data)
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

	postgresConn, err := postgres.Acquire(ctx)
	if err != nil {
		return nil, err
	}

	bus := redivent.NewBus(redisClient, jsoncodec.New[event.Event]())
	depsQueue := redipool.New[string](redisClient, strcodec.Codec{}, "deps_queue")
	deps := postgredeps.New(postgresConn, depsQueue)
	tasks := redikv.New[kantoku.Task](redisClient, jsoncodec.New[kantoku.Task](), "tasks")
	cells := redicell.NewStorage[[]byte](redisClient, bincodec.Codec{})

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

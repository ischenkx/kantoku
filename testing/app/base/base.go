package base

import (
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
	"kantoku"
	"kantoku/core/event"
	"kantoku/framework/depot"
	"kantoku/impl/common/codec/bincodec"
	"kantoku/impl/common/codec/jsoncodec"
	"kantoku/impl/common/codec/strcodec"
	redikv "kantoku/impl/common/data/kv/redis"
	"kantoku/impl/common/deps/postgredeps"
	redipool "kantoku/impl/common/pool/redis"
	redicell "kantoku/impl/core/cell/redis"
	redivent "kantoku/impl/core/event/redis"
)

type Base struct {
	Redis    *redis.Client
	Postgres *pgxpool.Pool
	Deps     *postgredeps.Deps
	Depot    *depot.Depot
	Kantoku  *kantoku.Kantoku
}

func redisClient() *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})
}

func PostgreClient(ctx context.Context) (*pgxpool.Pool, error) {
	return pgxpool.New(ctx, "postgres://postgres:root@localhost:5432/")
}

func New(ctx context.Context) (Base, error) {
	r := redisClient()
	p, err := PostgreClient(ctx)
	if err != nil {
		return Base{}, err
	}
	depsQueue := redipool.New[string](r, strcodec.Codec{}, "deps")
	deps := postgredeps.New(p, depsQueue)

	group2task := redikv.New[string](r, strcodec.Codec{}, "group2task")

	depotClient := depot.New(deps, group2task)

	kan := kantoku.Builder{
		Tasks: redikv.New[kantoku.StoredTask](
			r,
			jsoncodec.New[kantoku.StoredTask](),
			"tasks",
		),
		Cells:     redicell.New[[]byte](r, bincodec.Codec{}),
		Events:    redivent.NewBus(r, jsoncodec.New[event.Event]()),
		Scheduler: depotClient,
	}.Build()

	return Base{
		Redis:    r,
		Postgres: p,
		Deps:     deps,
		Depot:    depotClient,
		Kantoku:  kan,
	}, nil
}

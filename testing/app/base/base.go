package base

import (
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
	"kantoku"
	"kantoku/framework/depot"
	delay2 "kantoku/framework/depot/delay"
	subtask2 "kantoku/framework/depot/taskdep"
	"kantoku/framework/output"
	"kantoku/framework/status"
	"kantoku/impl/common/codec/bincodec"
	"kantoku/impl/common/codec/jsoncodec"
	"kantoku/impl/common/codec/strcodec"
	rebimap "kantoku/impl/common/data/bimap/redis"
	"kantoku/impl/common/data/cron/simple"
	redikv "kantoku/impl/common/data/kv/redis"
	"kantoku/impl/common/data/pool/redis"
	"kantoku/impl/common/deps/postgredeps"
	redicell "kantoku/impl/core/cell/redis"
	redivent "kantoku/impl/core/event/redis"
	"kantoku/platform"
	"kantoku/platform/event"
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

	groupTaskBimap := rebimap.NewBimap[string, string](
		"group2task",
		"task2group",
		strcodec.Codec{},
		strcodec.Codec{},
		r,
	)

	depotClient := depot.New(deps, groupTaskBimap)

	kan := kantoku.Builder{
		Tasks: redikv.New[kantoku.TaskInstance](
			r,
			jsoncodec.New[kantoku.TaskInstance](),
			"tasks",
		),
		Cells:     redicell.New[[]byte](r, bincodec.Codec{}),
		Events:    redivent.NewBus(r, jsoncodec.New[event.Event]()),
		Scheduler: depotClient,
	}.Build()

	outputs := redikv.New[platform.Result](r, jsoncodec.New[platform.Result](), "outputs")
	statusDB := redikv.New[status.Status](r, jsoncodec.New[status.Status](), "statuses")

	subtaskManager := subtask2.NewManager(deps, redikv.New[string](r, strcodec.Codec{}, "subtasks"))

	cronInputs := redipool.New[simple.Event](r, jsoncodec.New[simple.Event](), "cron_inputs")
	cronOutputs := redipool.New[string](r, strcodec.Codec{}, "cron_outputs")
	redisCron := simple.NewClient(cronInputs, cronOutputs)
	delayManager := delay2.NewManager(redisCron, deps)

	kan.Register(delay2.NewPlugin(delayManager))
	kan.Register(subtask2.NewPlugin(subtaskManager))
	kan.Register(output.NewPlugin(outputs))
	kan.Register(status.NewPlugin(statusDB))
	kan.Register(depot.NewPlugin(depotClient))

	return Base{
		Redis:    r,
		Postgres: p,
		Deps:     deps,
		Depot:    depotClient,
		Kantoku:  kan,
	}, nil
}

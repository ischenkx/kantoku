package common

import (
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
	"kantoku"
	"kantoku/impl/common/codec/jsoncodec"
	"kantoku/impl/common/codec/strcodec"
	rebimap "kantoku/impl/common/data/bimap/redis"
	redikv "kantoku/impl/common/data/kv/redis"
	redipool "kantoku/impl/common/data/pool/redis"
	"kantoku/impl/deps/postgredeps"
	redivent "kantoku/impl/platform/event/redis"
	"kantoku/platform"
	"kantoku/unused/backend/framework/depot"
	"log"
)

var redisClient redis.UniversalClient = nil

func makeRedisClient() redis.UniversalClient {
	if redisClient != nil {
		return redisClient
	}
	client := redis.NewClient(&redis.Options{
		Addr:     "redis:6379", // Redis server address
		Password: "",           // Redis server password (leave empty if not set)
		DB:       0,            // Redis database index
	})

	if cmd := client.Ping(context.Background()); cmd.Err() != nil {
		panic("failed to ping the redis client: " + cmd.Err().Error())
	}

	redisClient = client

	return client
}

func makePostgresClient(ctx context.Context) *pgxpool.Pool {
	client, err := pgxpool.New(ctx, "postgres://postgres:51413@postgres:5432/")

	if err != nil {
		panic("failed to create postgres deps: " + err.Error())
	}

	if err := client.Ping(ctx); err != nil {
		panic("failed to make ping postgres: " + err.Error())
	}

	return client
}

func MakeDeps() *postgredeps.Deps {
	pg := makePostgresClient(context.Background())
	deps := postgredeps.New(
		pg,
		redipool.New[string](makeRedisClient(), strcodec.Codec{}, "depot_groups"),
	)

	if err := deps.InitTables(context.Background()); err != nil {
		log.Println("failed to init tables:", err)
	}

	return deps
}

func MakeInputs() *depot.Depot {
	return depot.New(
		MakeDeps(),
		rebimap.NewBimap[string, string](
			"keys___",
			"values___",
			strcodec.Codec{},
			strcodec.Codec{},
			makeRedisClient(),
		),
		redipool.New[string](makeRedisClient(), strcodec.Codec{}, "TEST_STAND_INPUTS"),
	)
}

func makeOutputs() platform.Outputs {
	return redikv.New[platform.Result](makeRedisClient(), jsoncodec.New[platform.Result](), "TEST_STAND_OUTPUTS")
}

func makeBroker() platform.Broker {
	return redivent.New(jsoncodec.New[platform.Event](), makeRedisClient())
}

func makeDB() platform.DB[kantoku.TaskInstance] {
	return redikv.New[kantoku.TaskInstance](makeRedisClient(), jsoncodec.New[kantoku.TaskInstance](), "TEST_STAND_TASKS_DB")
}

func MakePlatform() platform.Platform[kantoku.TaskInstance] {
	return platform.New[kantoku.TaskInstance](
		makeDB(),
		MakeInputs(),
		makeOutputs(),
		makeBroker(),
	)
}

package common

import (
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
	"kantoku"
	"kantoku/common/data/bimap"
	"kantoku/common/data/kv"
	"kantoku/common/data/pool"
	"kantoku/framework/future"
	"kantoku/framework/plugins/depot"
	"kantoku/impl/common/codec/jsoncodec"
	"kantoku/impl/common/codec/strcodec"
	rebimap "kantoku/impl/common/data/bimap/redis"
	redikv "kantoku/impl/common/data/kv/redis"
	redipool "kantoku/impl/common/data/pool/redis"
	"kantoku/impl/deps/postgres/instant"
	redivent "kantoku/impl/platform/event/redis"
	redismeta "kantoku/impl/plugins/meta/redis"
	"kantoku/kernel"
	"kantoku/kernel/platform"
	"log"
)

type futureRunner struct {
	queue pool.Writer[future.ID]
}

func (f futureRunner) Run(ctx context.Context, resolution future.Resolution) {
	if err := f.queue.Write(ctx, resolution.Future.ID); err != nil {
		log.Println("failed to put a resolution in the queue:", err)
	}
}

var redisClient redis.UniversalClient = nil

func MakeRedisClient() redis.UniversalClient {
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

func MakePostgresClient(ctx context.Context) *pgxpool.Pool {
	client, err := pgxpool.New(ctx, "postgres://postgres:51413@postgres:5432/")

	if err != nil {
		panic("failed to create postgres deps: " + err.Error())
	}

	if err := client.Ping(ctx); err != nil {
		panic("failed to make ping postgres: " + err.Error())
	}

	return client
}

func MakeDeps() *instant.Deps {
	pg := MakePostgresClient(context.Background())
	deps := instant.New(
		pg,
		redipool.New[string](MakeRedisClient(), strcodec.Codec{}, "depot_groups"),
	)

	if err := deps.InitTables(context.Background()); err != nil {
		log.Println("failed to init tables:", err)
	}

	return deps
}

func MakeDepotBimap() bimap.Bimap[string, string] {
	return rebimap.NewBimap[string, string](
		"keys___",
		"values___",
		strcodec.Codec{},
		strcodec.Codec{},
		MakeRedisClient(),
	)
}

func MakeInputs() *depot.Depot {
	return depot.New(
		MakeDeps(),
		MakeDepotBimap(),
		redipool.New[string](MakeRedisClient(), strcodec.Codec{}, "TEST_STAND_INPUTS"),
	)
}

func MakeOutputs() platform.Outputs {
	return redikv.New[platform.Result](MakeRedisClient(), jsoncodec.New[platform.Result](), "TEST_STAND_OUTPUTS")
}

func MakeBroker() platform.Broker {
	return redivent.New(jsoncodec.New[platform.Event](), MakeRedisClient())
}

func MakeDB() platform.DB[kernel.Task] {
	return redikv.New[kernel.Task](MakeRedisClient(), jsoncodec.New[kernel.Task](), "TEST_STAND_TASKS_DB")
}

func MakePlatform() platform.Platform[kernel.Task] {
	return platform.New[kernel.Task](
		MakeDB(),
		MakeInputs(),
		MakeOutputs(),
		MakeBroker(),
	)
}

func MakeFutureResolutionQueue() pool.Pool[future.ID] {
	return redipool.New[future.ID](MakeRedisClient(), strcodec.Codec{}, "TEST_STAND_FUTURE_RESOLUTIONS_QUEUE")
}

func MakeFuturesManager() *future.Manager {
	return future.NewManager(
		redikv.New[future.Future](MakeRedisClient(), jsoncodec.New[future.Future](), "TEST_STAND_FUTURES"),
		redikv.New[future.Resource](MakeRedisClient(), jsoncodec.New[future.Resource](), "TEST_STAND_FUTURE_RESOLUTIONS"),
		futureRunner{queue: MakeFutureResolutionQueue()},
	)
}

func MakeTaskDepDB() kv.Database[string, string] {
	return redikv.New[string](
		MakeRedisClient(),
		strcodec.Codec{},
		"TEST_STAND_TASK_DEPS",
	)
}

func MakeFutDepDB() kv.Database[string, string] {
	return redikv.New[string](
		MakeRedisClient(),
		strcodec.Codec{},
		"TEST_STAND_FUT_DEPS",
	)
}

func MakeKantoku() *kantoku.Kantoku {
	return kantoku.NewBuilder().
		ConfigureParametrizationCodec(jsoncodec.New[kantoku.Parametrization]()).
		ConfigureSettings(
			kantoku.Settings{AutoInputDependencies: true},
		).
		ConfigurePlatform(MakePlatform()).
		//ConfigureContexts().
		ConfigureFutures(MakeFuturesManager()).
		ConfigureTaskdep(MakeTaskDepDB()).
		ConfigureFutdep(MakeFutDepDB()).
		ConfigureDepot(MakeDepotBimap()).
		ConfigureDeps(MakeDeps()).
		ConfigureMeta(redismeta.NewDB("META", MakeRedisClient()), jsoncodec.Dynamic{}).
		Build()
}

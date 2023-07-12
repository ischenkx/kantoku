package common

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"kantoku"
	"kantoku/backend/executor/evexec"
	"kantoku/common/data/bimap"
	"kantoku/common/data/kv"
	"kantoku/common/data/pool"
	"kantoku/common/data/record"
	"kantoku/framework/future"
	"kantoku/framework/infra/demon"
	"kantoku/framework/plugins/info"
	"kantoku/framework/utils/demons"
	"kantoku/impl/common/codec/jsoncodec"
	"kantoku/impl/common/codec/strcodec"
	redimap "kantoku/impl/common/data/bimap/redis"
	redikv "kantoku/impl/common/data/kv/redis"
	redipool "kantoku/impl/common/data/pool/redis"
	mongorec "kantoku/impl/common/data/record/mongo"
	"kantoku/impl/deps/postgres/instant"
	redivent "kantoku/impl/platform/event/redis"
	"kantoku/kernel"
	"kantoku/kernel/platform"
	"log"
	"strconv"
)

var redisClient redis.UniversalClient = nil

func MakeRedisClient() redis.UniversalClient {
	if redisClient != nil {
		return redisClient
	}
	client := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379", // Redis server address
		Password: "",               // Redis server password (leave empty if not set)
		DB:       0,                // Redis database index
	})

	if cmd := client.Ping(context.Background()); cmd.Err() != nil {
		panic("failed to ping the redis client: " + cmd.Err().Error())
	}

	redisClient = client

	return client
}

func MakePostgresClient(ctx context.Context) *pgxpool.Pool {
	client, err := pgxpool.New(ctx, "postgres://postgres:51413@localhost:5432/")

	if err != nil {
		panic("failed to create postgres deps: " + err.Error())
	}

	if err := client.Ping(ctx); err != nil {
		panic("failed to make ping postgres: " + err.Error())
	}

	return client
}

func MakeMongoClient(ctx context.Context) *mongo.Client {
	// Set connection configurations
	clientOptions := options.Client().ApplyURI("mongodb://localhost:27017")

	// Connect to the MongoDB server
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		log.Fatal(err)
	}
	go func() {
		<-ctx.Done()

		if err = client.Disconnect(ctx); err != nil {
			log.Println("failed to disconnect from mongodb:", err)
		}
	}()

	if err := client.Ping(ctx, readpref.Nearest()); err != nil {
		log.Fatal(err)
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
	return redimap.NewBimap[string, string](
		"keys___",
		"values___",
		strcodec.Codec{},
		strcodec.Codec{},
		MakeRedisClient(),
	)
}

func MakeInputsQueue() pool.Pool[string] {
	return redipool.New[string](MakeRedisClient(), strcodec.Codec{}, "TEST_STAND_INPUTS")
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
		MakeInputsQueue(),
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
		MakeFutureResolutionQueue(),
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

func MakeInfoRecords() record.Storage {
	client := MakeMongoClient(context.Background())
	database := client.Database("test_stand")
	collection := database.Collection("info_records")

	return mongorec.New(collection)
}

func MakeKantoku() (*kantoku.Kantoku, error) {
	kan := &kantoku.Kantoku{}
	kan1, err := kantoku.Configure().
		Parametrization(jsoncodec.New[kantoku.Parametrization]()).
		Settings(kantoku.Settings{AutoInputDependencies: true}).
		Platform(MakePlatform()).
		Futures(MakeFuturesManager()).
		TaskDep(MakeTaskDepDB()).
		FutDep(MakeFutDepDB()).
		Depot(MakeDepotBimap()).
		Deps(MakeDeps()).
		Deployer(deployer{}).
		Demons(
			demons.Functional("EXECUTOR",
				func(ctx context.Context) error {
					return evexec.Builder[kernel.Task]{
						Runner:   runner{kantoku: kan},
						Platform: kan.Kernel().Platform(),
						Resolver: evexec.ConstantResolver("TEST_STAND_TASK_EVENTS"),
					}.Build().Run(ctx)
				},
			),
		).
		Info(MakeInfoRecords(), info.Settings{IdProperty: "task_id"}).
		Compile()
	if err != nil {
		return nil, fmt.Errorf("failed to build kantoku: %s", err)
	}
	*kan = *kan1

	return kan, nil
}

type deployer struct{}

func (r deployer) Deploy(ctx context.Context, demons ...demon.Demon) error {
	for _, dem := range demons {
		name := dem.Name
		log.Println("Starting a demon:", name)
		if fn, ok := dem.Parameter.(func(context.Context) error); ok {
			go func() {
				if err := fn(ctx); err != nil {
					log.Printf("[%s] failed to run: %s\n", name, err)
				}
			}()
			log.Println("Started!")
		} else {
			log.Println("Parameter is not a function... Can't run it")
		}
		log.Println("---------")
	}
	<-ctx.Done()
	return nil
}

type runner struct {
	kantoku *kantoku.Kantoku
}

func (r runner) Run(ctx context.Context, raw kernel.Task) ([]byte, error) {
	task := r.kantoku.Task(raw.ID())
	inputs, err := task.Inputs(ctx)
	if err != nil {
		log.Println("failed to load inputs:", err)
		return nil, err
	}

	outputs, err := task.Outputs(ctx)
	if err != nil {
		log.Println("failed to load outputs:", err)
		return nil, err
	}

	switch raw.Type {
	case "factorial":
		input, err := r.kantoku.Futures().Load(ctx, inputs[0])
		if err != nil {
			log.Println("failed to load input:", err)
			return nil, err
		}

		num, err := strconv.Atoi(string(input.Resource))
		if err != nil {
			log.Println("failed to cast input to a number:", err)
			return nil, err
		}

		res := fact(num)
		err = r.kantoku.Futures().Resolve(ctx, outputs[0], []byte(strconv.Itoa(res)))
		if err != nil {
			log.Println("failed to resolve a given output future:", err)
			return nil, err
		}

		log.Println("TASK:", raw.Type)
		log.Println("Input:", num)
		log.Println("Result:", res)
	case "mul":
		input1, err := r.kantoku.Futures().Load(ctx, inputs[0])
		if err != nil {
			log.Println("failed to load input1:", err)
			return nil, err
		}

		input2, err := r.kantoku.Futures().Load(ctx, inputs[1])
		if err != nil {
			log.Println("failed to load input2:", err)
			return nil, err
		}

		num1, err := strconv.Atoi(string(input1.Resource))
		if err != nil {
			log.Println("failed to cast input1 to a number:", err)
			return nil, err
		}

		num2, err := strconv.Atoi(string(input2.Resource))
		if err != nil {
			log.Println("failed to cast input2 to a number:", err)
			return nil, err
		}

		res := num1 * num2
		err = r.kantoku.Futures().Resolve(ctx, outputs[0], []byte(strconv.Itoa(res)))
		if err != nil {
			log.Println("failed to resolve a given output future:", err)
			return nil, err
		}
		log.Println("TASK:", raw.Type)
		log.Println("Input1:", num1)
		log.Println("Input2:", num2)
		log.Println("Result:", res)
	}

	return nil, nil
}

func fact(x int) int {
	if x <= 1 {
		return 1
	}
	return x * fact(x-1)
}

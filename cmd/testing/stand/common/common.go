package common

import (
	"context"
	"fmt"
	"github.com/ThreeDotsLabs/watermill-nats/v2/pkg/nats"
	codec "github.com/ischenkx/kantoku/pkg/common/data/codec"
	"github.com/ischenkx/kantoku/pkg/core/event"
	"github.com/ischenkx/kantoku/pkg/core/resource"
	"github.com/ischenkx/kantoku/pkg/core/system"
	"github.com/ischenkx/kantoku/pkg/core/task"
	"github.com/ischenkx/kantoku/pkg/lib/impl/broker/watermill"
	redisResources "github.com/ischenkx/kantoku/pkg/lib/impl/core/resource/redis"
	mongorec "github.com/ischenkx/kantoku/pkg/lib/impl/data/record/mongo"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/kelseyhightower/envconfig"
	"github.com/lmittmann/tint"
	"github.com/redis/go-redis/v9"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log/slog"
	"os"
	"time"
)

type Config struct {
	RedisHost        string `envconfig:"REDIS_HOST" default:"localhost"`
	RedisPort        int    `envconfig:"REDIS_PORT" default:"6379"`
	MongoHost        string `envconfig:"MONGO_HOST" default:"localhost"`
	MongoPort        int    `envconfig:"MONGO_PORT" default:"27017"`
	PostgresHost     string `envconfig:"POSTGRES_HOST" default:"localhost"`
	PostgresPort     int    `envconfig:"POSTGRES_PORT" default:"5432"`
	PostgresUser     string `envconfig:"POSTGRES_USER" default:"postgres"`
	PostgresPassword string `envconfig:"POSTGRES_PASSWORD" default:"postgres"`
	NatsURL          string `envconfig:"NATS_URL" default:"nats://localhost:4222"`
}

func InitLogger() {
	slog.SetDefault(
		slog.New(tint.NewHandler(os.Stdout, &tint.Options{
			Level:      slog.LevelDebug,
			TimeFormat: time.Kitchen,
		})),
	)
}

func NewRedis(ctx context.Context, host string, port int) redis.UniversalClient {
	client := redis.NewUniversalClient(&redis.UniversalOptions{
		Addrs:       []string{fmt.Sprintf("%s:%d", host, port)},
		DialTimeout: time.Minute,
	})

	if err := client.Ping(ctx).Err(); err != nil {
		panic(err)
	}

	return client
}

func NewMongo(ctx context.Context, host string, port int) *mongo.Client {
	url := fmt.Sprintf("mongodb://%s:%d", host, port)
	clientOptions := options.Client().ApplyURI(url)

	// Connect to the MongoDB server
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		panic(err)
	}

	return client
}

func NewPostgres(ctx context.Context, host string, port int, user, pwd string) *pgxpool.Pool {
	pool, err := pgxpool.New(ctx, fmt.Sprintf("postgres://%s:%s@%s:%d/", user, pwd, host, port))
	if err != nil {
		panic(fmt.Sprintf("failed to create pool: %s", err))
	}

	return pool
}

func NewConfig() Config {
	var config Config
	err := envconfig.Process("", &config)
	if err != nil {
		panic(err)
	}

	return config
}

func NewSystem(ctx context.Context, consumer string) *system.System {
	config := NewConfig()

	fmt.Printf("Connecting (mongo=%s:%d redis=%s:%d nats=%s)\n",
		config.MongoHost,
		config.MongoPort,
		config.RedisHost,
		config.RedisPort,
		config.NatsURL)

	mongoClient := NewMongo(ctx, config.MongoHost, config.MongoPort)
	redisClient := NewRedis(ctx, config.RedisHost, config.RedisPort)

	//brokerAgent, err := watermill.Redis(
	//	redisClient,
	//	redisstream.SubscriberConfig{
	//		// TODO move to the constructor
	//		Consumer: consumer,
	//	},
	//	redisstream.PublisherConfig{},
	//)
	//if err != nil {
	//	panic(fmt.Sprintf("failed to create a broker agent: %w", err))
	//}

	brokerAgent, err := watermill.Nats(
		config.NatsURL,
		nats.SubscriberConfig{},
		nats.PublisherConfig{},
	)
	if err != nil {
		panic(fmt.Sprintf("failed to create a broker agent: %w", err))
	}

	broker := watermill.Broker[event.Event]{
		Agent:                     brokerAgent,
		ItemCodec:                 codec.JSON[event.Event](),
		ConsumerChannelBufferSize: 1024,
	}

	return system.New(
		event.NewBroker(broker),
		redisResources.New(redisClient, codec.JSON[resource.Resource](), "test-resources"),
		mongorec.New[task.Task](
			mongoClient.
				Database("testing").
				Collection("task_info"),
			task.Codec{},
		),
	)
}

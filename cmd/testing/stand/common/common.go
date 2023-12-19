package common

import (
	"context"
	codec "github.com/ischenkx/kantoku/pkg/common/data/codec"
	mongorec "github.com/ischenkx/kantoku/pkg/impl/data/record/mongo"
	redisEvents "github.com/ischenkx/kantoku/pkg/impl/kernel/event/redis"
	redisResources "github.com/ischenkx/kantoku/pkg/impl/kernel/resource/redis"
	redisTasks "github.com/ischenkx/kantoku/pkg/impl/kernel/task/redis"
	"github.com/ischenkx/kantoku/pkg/system"
	"github.com/ischenkx/kantoku/pkg/system/kernel/resource"
	"github.com/ischenkx/kantoku/pkg/system/kernel/task"
	"github.com/lmittmann/tint"
	"github.com/redis/go-redis/v9"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log/slog"
	"os"
	"time"
)

func InitLogger() {
	slog.SetDefault(
		slog.New(tint.NewHandler(os.Stdout, &tint.Options{
			Level:      slog.LevelDebug,
			TimeFormat: time.Kitchen,
		})),
	)
}

func NewRedis(ctx context.Context) redis.UniversalClient {
	client := redis.NewUniversalClient(&redis.UniversalOptions{
		Addrs: []string{"172.23.146.206:6379"},
		//Addrs: []string{":6379"},
	})

	if err := client.Ping(ctx).Err(); err != nil {
		panic(err)
	}

	return client
}

func NewMongo(ctx context.Context) *mongo.Client {
	clientOptions := options.Client().ApplyURI("mongodb://localhost:27017")

	// Connect to the MongoDB server
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		panic(err)
	}

	return client
}

func NewSystem(ctx context.Context, consumer string) *system.System {
	mongoClient := NewMongo(ctx)
	redisClient := NewRedis(ctx)

	return system.New(
		redisEvents.New(redisClient, redisEvents.StreamSettings{
			BatchSize:         64,
			ChannelBufferSize: 1024,
			Consumer:          consumer,
		}),
		redisResources.New(redisClient, codec.JSON[resource.Resource](), "test-resources"),
		redisTasks.New(redisClient, codec.JSON[task.Task](), "test-tasks"),
		mongorec.New(
			mongoClient.
				Database("testing").
				Collection("task_info"),
		),
	)
}

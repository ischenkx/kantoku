package main

import (
	"context"
	"github.com/ischenkx/kantoku/cmd/testing/stand/common"
	"github.com/ischenkx/kantoku/pkg/impl/data/dependency/postgres/batched"
	"github.com/ischenkx/kantoku/pkg/processors/scheduler/dependencies/simple"
	"github.com/ischenkx/kantoku/pkg/processors/scheduler/dependencies/simple/manager"
	resourceResolver "github.com/ischenkx/kantoku/pkg/processors/scheduler/dependencies/simple/manager/resolvers/resource_resolver"
	"github.com/ischenkx/kantoku/pkg/processors/scheduler/dependencies/simple/manager/task2group"
	"log/slog"
	"time"
)

//
//func main() {
//	common.InitLogger()
//
//	slog.Info("Starting...")
//	err := dummy.NewProcessor(common.NewSystem(context.Background(), "scheduler-0")).Process(context.Background())
//	if err != nil {
//		slog.Error("failed", slog.String("error", err.Error()))
//	}
//}

func main() {
	common.InitLogger()
	ctx := context.Background()

	config := common.NewConfig()

	dependencies := batched.New(
		common.NewPostgres(ctx,
			config.PostgresHost,
			config.PostgresPort,
			config.PostgresUser,
			config.PostgresPassword),
		batched.Config{
			PollingInterval:  time.Second,
			PollingBatchSize: 256,
		})

	system := common.NewSystem(context.Background(), "scheduler-0")

	mongo := common.NewMongo(ctx, config.MongoHost, config.MongoPort)

	processor := simple.NewProcessor(
		system,
		dependencies,
		&task2group.RedisStorage{
			Client: common.NewRedis(
				ctx,
				config.RedisHost,
				config.RedisPort,
			),
		},
		map[string]manager.Resolver{
			"resource": &resourceResolver.Resolver{
				System: system,
				Storage: &resourceResolver.MongoStorage{
					Collection:  mongo.Database("testing").Collection("resource_dependencies"),
					PollTimeout: time.Minute * 5,
				},
				PollLimit:    1024,
				PollInterval: time.Second,
			},
		},
	)

	slog.Info("Starting...")
	err := processor.Process(context.Background())
	if err != nil {
		slog.Error("failed", slog.String("error", err.Error()))
	}
}

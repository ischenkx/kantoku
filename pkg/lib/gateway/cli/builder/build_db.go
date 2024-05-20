package builder

import (
	"context"
	"fmt"
	"github.com/ischenkx/kantoku/pkg/lib/gateway/cli/config"
	"github.com/ischenkx/kantoku/pkg/lib/gateway/cli/errx"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoInfo struct {
	Client     *mongo.Client
	Collection string
	Database   string
}

func (info MongoInfo) GetCollection() *mongo.Collection {
	return info.Client.Database(info.Database).Collection(info.Collection)
}

func (builder *Builder) BuildPostgres(ctx context.Context, cfg config.DynamicConfig) (*pgxpool.Pool, error) {
	var connConfig struct {
		URI string
	}
	if err := cfg.Bind(&connConfig); err != nil {
		return nil, errx.FailedToBind(err)
	}

	pool, err := pgxpool.New(context.Background(), connConfig.URI)
	if err != nil {
		return nil, fmt.Errorf("failed to create a pool: %w", err)
	}

	return pool, nil
}

func (builder *Builder) BuildRedis(ctx context.Context, cfg config.DynamicConfig) (redis.UniversalClient, error) {
	var connConfig struct {
		URI string
	}
	if err := cfg.Bind(&connConfig); err != nil {
		return nil, errx.FailedToBind(err)
	}

	client := redis.NewUniversalClient(&redis.UniversalOptions{
		Addrs: []string{
			connConfig.URI,
		},
	})

	return client, nil
}

func (builder *Builder) BuildMongo(ctx context.Context, cfg config.DynamicConfig) (MongoInfo, error) {
	var connConfig struct {
		URI        string
		DB         string
		Collection string
	}
	if err := cfg.Bind(&connConfig); err != nil {
		return MongoInfo{}, errx.FailedToBind(err)
	}

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(connConfig.URI))
	if err != nil {
		return MongoInfo{}, fmt.Errorf("failed to connect to mongo: %w", err)
	}

	return MongoInfo{
		Client:     client,
		Collection: connConfig.Collection,
		Database:   connConfig.DB,
	}, nil
}

//func (builder *Builder)BuildNats(ctx context.Context, cfg config.DynamicConfig) {
//
//}

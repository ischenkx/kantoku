package resourceResolverV2

import (
	"context"
	"errors"
	"fmt"
	"github.com/ischenkx/kantoku/pkg/common/data"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"time"
)

type MongoStorage struct {
	Collection *mongo.Collection
}

func (storage *MongoStorage) Save(ctx context.Context, dependencyId string, resourceId string) error {
	_, err := storage.Collection.InsertOne(ctx, bson.M{
		"dependency_id": dependencyId,
		"resource_id":   resourceId,
		"status":        "pending",
		"updated_at":    time.Now().Unix(),
	})
	return err
}

func (storage *MongoStorage) Load(ctx context.Context, resourceId string) (dependencyId string, err error) {
	result := storage.Collection.FindOne(ctx, bson.M{"resource_id": resourceId})
	err = result.Err()
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return "", data.NotFoundErr
		}
		return "", fmt.Errorf("failed to find: %w", err)
	}

	var raw bson.M
	if err := result.Decode(&raw); err != nil {
		return "", fmt.Errorf("failed to decode: %w", err)
	}

	rawDependencyId, ok := raw["dependency_id"]
	if !ok {
		return "", fmt.Errorf("no dependency_id in doc")
	}

	dependencyId, ok = rawDependencyId.(string)
	if !ok {
		return "", fmt.Errorf("dependency_id is not a string")
	}

	return dependencyId, nil
}

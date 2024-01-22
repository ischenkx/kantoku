package resourceResolver

import (
	"context"
	"fmt"
	"github.com/samber/lo"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

type MongoStorage struct {
	Collection  *mongo.Collection
	PollTimeout time.Duration
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

func (storage *MongoStorage) Resolve(ctx context.Context, dependencyIds ...string) error {
	if len(dependencyIds) == 0 {
		return nil
	}
	_, err := storage.Collection.UpdateMany(ctx,
		bson.M{
			"status": bson.M{
				"$in": bson.A{"pending", "polled"},
			},
			"dependency_id": bson.M{
				"$in": dependencyIds,
			},
		},
		bson.M{
			"$set": bson.M{
				"status":     "resolved",
				"updated_at": time.Now().Unix(),
			},
		})

	return err
}

func (storage *MongoStorage) Poll(ctx context.Context, limit int) ([]Binding, error) {
	updatedAtLowerBoundary := time.Now().Add(-storage.PollTimeout).Unix()

	cursor, err := storage.Collection.Find(
		ctx,
		bson.M{
			"status": "pending",
			"updated_at": bson.M{
				"$gte": updatedAtLowerBoundary,
			},
		},
		options.Find().SetLimit(int64(limit)))
	if err != nil {
		return nil, fmt.Errorf("failed to find pending bindings: %w", err)
	}
	defer cursor.Close(ctx)

	var rawBindings []bson.M
	if err := cursor.All(ctx, &rawBindings); err != nil {
		return nil, fmt.Errorf("failed to decode bindings: %w", err)
	}

	bindings := lo.Map(rawBindings, func(rawBinding bson.M, _ int) Binding {
		return Binding{
			DependencyId: rawBinding["dependency_id"].(string),
			ResourceId:   rawBinding["resource_id"].(string),
		}
	})

	return bindings, nil
}

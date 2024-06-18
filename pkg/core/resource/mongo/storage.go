package mongo

import (
	"context"
	"fmt"
	"github.com/ischenkx/kantoku/pkg/common/data/uid"
	"github.com/ischenkx/kantoku/pkg/core/resource"
	"github.com/samber/lo"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"time"
)

var _ resource.Storage = (*Storage)(nil)

type Storage struct {
	Collection  *mongo.Collection
	LockTimeout time.Duration
}

func (storage *Storage) Load(ctx context.Context, ids ...resource.ID) ([]resource.Resource, error) {
	cursor, err := storage.Collection.Find(ctx,
		bson.M{
			"_id": bson.M{
				"$in": ids,
			},
		})
	if err != nil {
		return nil, fmt.Errorf("failed to find: %w", err)
	}

	var resources []Resource
	if err := cursor.All(ctx, &resources); err != nil {
		return nil, fmt.Errorf("failed to decode: %w", err)
	}

	index := lo.SliceToMap(resources, func(res Resource) (string, resource.Resource) {
		return res.ID, resource.Resource{
			Data:   res.Data,
			ID:     res.ID,
			Status: resource.Status(res.Status),
		}
	})

	result := lo.Map(ids, func(id string, _ int) resource.Resource {
		if res, ok := index[id]; ok {
			return res
		}

		return resource.Resource{
			ID:     id,
			Status: resource.DoesNotExist,
		}
	})

	return result, nil
}

func (storage *Storage) Alloc(ctx context.Context, amount int) ([]resource.ID, error) {
	ids := make([]string, 0, amount)

	resources := lo.RepeatBy(amount, func(_ int) any {
		id := uid.Generate()
		ids = append(ids, id)
		return Resource{
			ID:     id,
			Status: resource.Allocated,
		}
	})

	if _, err := storage.Collection.InsertMany(ctx, resources); err != nil {
		return nil, fmt.Errorf("failed to insert: %w", err)
	}

	return ids, nil
}

func (storage *Storage) Init(ctx context.Context, resources []resource.Resource) error {
	resources = lo.UniqBy(resources, func(res resource.Resource) string {
		return res.ID
	})

	// Generate a new version
	version := uid.Generate()
	now := time.Now()

	// Lock resources and load
	storage.Collection.UpdateMany(
		ctx,
		bson.M{},
		bson.M{
			"version":    version,
			"updated_at": now,
		},
	)

	storage.Collection

	// Validate resources
	// Update resources
	// Unlock resources
}

func (storage *Storage) Dealloc(ctx context.Context, ids []resource.ID) error {
	_, err := storage.Collection.DeleteMany(ctx, bson.M{
		"_id": bson.M{
			"$in": ids,
		},
	})
	if err != nil {
		return fmt.Errorf("failed to delete: %w", err)
	}

	return nil
}

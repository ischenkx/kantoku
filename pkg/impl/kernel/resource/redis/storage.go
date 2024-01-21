package redis

import (
	"context"
	"errors"
	"fmt"
	"github.com/google/uuid"
	codec "github.com/ischenkx/kantoku/pkg/common/data/codec"
	"github.com/ischenkx/kantoku/pkg/system/kernel/resource"
	"github.com/redis/go-redis/v9"
	"github.com/samber/lo"
	"strings"
)

type Storage struct {
	client    redis.UniversalClient
	codec     codec.Codec[resource.Resource, []byte]
	keyPrefix string
}

func New(client redis.UniversalClient, codec codec.Codec[resource.Resource, []byte], keyPrefix string) *Storage {
	return &Storage{
		client:    client,
		codec:     codec,
		keyPrefix: keyPrefix,
	}
}

func (storage *Storage) Load(ctx context.Context, ids ...resource.ID) ([]resource.Resource, error) {
	if len(ids) == 0 {
		return []resource.Resource{}, nil
	}

	cmd := storage.client.MGet(ctx, lo.Map(ids, func(id string, _ int) string {
		return storage.globalResourceID(id)
	})...)
	if cmd.Err() != nil {
		return nil, cmd.Err()
	}

	resources := make([]resource.Resource, 0, len(ids))

	for index, value := range cmd.Val() {
		id := ids[index]

		res, err := storage.parseResourceFromRedisValue(id, value)
		if err != nil {
			return nil, fmt.Errorf("failed to parse the resource (id='%s'): %w", id, err)
		}

		resources = append(resources, res)
	}

	return resources, nil
}

func (storage *Storage) Alloc(ctx context.Context, amount int) ([]resource.ID, error) {
	ids := lo.Times(amount, func(_ int) string {
		return storage.generateKey()
	})

	var arguments []any

	for _, id := range ids {
		encodedResource, err := storage.codec.Encode(resource.Resource{
			ID:     id,
			Status: resource.Allocated,
		})
		if err != nil {
			return nil, fmt.Errorf("failed to encode the resource: %w", err)
		}

		arguments = append(arguments, storage.globalResourceID(id), encodedResource)
	}

	cmd := storage.client.MSet(ctx, arguments)
	if cmd.Err() != nil {
		return nil, cmd.Err()
	}

	return ids, nil
}

func (storage *Storage) Init(ctx context.Context, resources []resource.Resource) error {
	resourceIDs := lo.Map(resources, func(res resource.Resource, _ int) string {
		return res.ID
	})
	globalResourceIDs := lo.Map(resourceIDs, func(localID string, _ int) string {
		return storage.globalResourceID(localID)
	})

	tx := func(tx *redis.Tx) error {
		loadedResources, err := tx.MGet(ctx, globalResourceIDs...).Result()
		if err != nil {
			return fmt.Errorf("failed to load resources: %w", err)
		}

		newValues := make([][]byte, 0, len(globalResourceIDs))

		for _, pair := range lo.Zip2(resources, loadedResources) {
			providedResource, rawLoadedResource := pair.A, pair.B
			loadedResource, err := storage.parseResourceFromRedisValue(providedResource.ID, rawLoadedResource)
			if err != nil {
				return fmt.Errorf("failed to parse the redis value (id='%s'): %w", providedResource.ID, err)
			}

			if loadedResource.Status != resource.Allocated {
				return fmt.Errorf("can't initialize a resource with status '%s' (id='%s')",
					loadedResource.Status,
					loadedResource.ID)
			}

			loadedResource.Data = providedResource.Data
			loadedResource.Status = resource.Ready

			encodedReadyResource, err := storage.codec.Encode(loadedResource)
			if err != nil {
				return fmt.Errorf("failed to encode a ready resource: %w", err)
			}

			newValues = append(newValues, encodedReadyResource)
		}

		_, err = tx.TxPipelined(ctx, func(pipeliner redis.Pipeliner) error {
			arguments := lo.Interleave[any](
				lo.ToAnySlice(globalResourceIDs),
				lo.ToAnySlice(newValues),
			)
			return pipeliner.MSet(ctx, arguments).Err()
		})

		return err
	}

	err := storage.client.Watch(ctx, tx, globalResourceIDs...)

	if err != nil {
		return fmt.Errorf("redis transaction failed: %w", err)
	}

	return nil
}

func (storage *Storage) Dealloc(ctx context.Context, ids []resource.ID) error {
	return storage.client.
		Del(ctx, lo.Map(ids, func(id string, _ int) string {
			return storage.globalResourceID(id)
		})...).
		Err()
}

func (storage *Storage) parseResourceFromRedisValue(id string, value any) (res resource.Resource, err error) {
	if value == nil {
		res = resource.Resource{
			ID:     resource.ID(id),
			Status: resource.DoesNotExist,
		}
	} else {
		encodedResource, ok := value.(string)
		if !ok {
			return res, errors.New("redis value is not a string")
		}

		res, err = storage.codec.Decode([]byte(encodedResource))
		if err != nil {
			return
		}
	}

	return
}

func (storage *Storage) generateKey() string {
	return strings.ReplaceAll(uuid.New().String(), "-", "")
}

func (storage *Storage) globalResourceID(key string) string {
	return fmt.Sprintf("%s:%s", storage.keyPrefix, key)
}

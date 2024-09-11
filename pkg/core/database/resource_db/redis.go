package resourcedb

import (
	"context"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/ischenkx/kantoku/pkg/common/data/codec"
	"github.com/ischenkx/kantoku/pkg/core"
	"github.com/redis/go-redis/v9"
	"github.com/samber/lo"
	"strings"
)

type RedisDB struct {
	client    redis.UniversalClient
	codec     codec.Codec[core.Resource, []byte]
	keyPrefix string
}

func NewRedisDB(client redis.UniversalClient, codec codec.Codec[core.Resource, []byte], keyPrefix string) *RedisDB {
	return &RedisDB{
		client:    client,
		codec:     codec,
		keyPrefix: keyPrefix,
	}
}

func (storage *RedisDB) Load(ctx context.Context, ids ...string) ([]core.Resource, error) {
	if len(ids) == 0 {
		return []core.Resource{}, nil
	}

	cmd := storage.client.MGet(ctx, lo.Map(ids, func(id string, _ int) string {
		return storage.globalResourceID(id)
	})...)
	if cmd.Err() != nil {
		return nil, cmd.Err()
	}

	resources := make([]core.Resource, 0, len(ids))

	for index, value := range cmd.Val() {
		id := ids[index]

		res, err := storage.parseResourceFromRedisValue(id, value)
		if err != nil {
			return nil, fmt.Errorf("failed to parse the resource_db (id='%s'): %w", id, err)
		}

		resources = append(resources, res)
	}

	return resources, nil
}

func (storage *RedisDB) Alloc(ctx context.Context, amount int) ([]string, error) {
	if amount == 0 {
		return []string{}, nil
	}

	ids := lo.Times(amount, func(_ int) string {
		return storage.generateKey()
	})

	var arguments []any

	for _, id := range ids {
		encodedResource, err := storage.codec.Encode(core.Resource{
			ID:     id,
			Status: core.ResourceStatuses.Allocated,
		})
		if err != nil {
			return nil, fmt.Errorf("failed to encode the resource_db: %w", err)
		}

		arguments = append(arguments, storage.globalResourceID(id), encodedResource)
	}

	cmd := storage.client.MSet(ctx, arguments)
	if cmd.Err() != nil {
		return nil, cmd.Err()
	}

	return ids, nil
}

func (storage *RedisDB) Init(ctx context.Context, resources []core.Resource) error {
	if len(resources) == 0 {
		return nil
	}
	resourceIDs := lo.Map(resources, func(res core.Resource, _ int) string {
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

			if loadedResource.Status != core.ResourceStatuses.Allocated {
				return fmt.Errorf("can't initialize a resource_db with status '%s' (id='%s')",
					loadedResource.Status,
					loadedResource.ID)
			}

			loadedResource.Data = providedResource.Data
			loadedResource.Status = core.ResourceStatuses.Ready

			encodedReadyResource, err := storage.codec.Encode(loadedResource)
			if err != nil {
				return fmt.Errorf("failed to encode a ready resource_db: %w", err)
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

func (storage *RedisDB) Dealloc(ctx context.Context, ids []string) error {
	if len(ids) == 0 {
		return nil
	}
	return storage.client.
		Del(ctx, lo.Map(ids, func(id string, _ int) string {
			return storage.globalResourceID(id)
		})...).
		Err()
}

func (storage *RedisDB) parseResourceFromRedisValue(id string, value any) (res core.Resource, err error) {
	if value == nil {
		res = core.Resource{
			ID:     id,
			Status: core.ResourceStatuses.DoesNotExist,
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

func (storage *RedisDB) generateKey() string {
	return strings.ReplaceAll(uuid.New().String(), "-", "")
}

func (storage *RedisDB) globalResourceID(key string) string {
	return fmt.Sprintf("%s:%s", storage.keyPrefix, key)
}

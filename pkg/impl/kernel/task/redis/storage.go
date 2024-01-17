package redis

import (
	"context"
	"errors"
	"fmt"
	"github.com/google/uuid"
	codec "github.com/ischenkx/kantoku/pkg/common/data/codec"
	"github.com/ischenkx/kantoku/pkg/system/kernel/task"
	"github.com/redis/go-redis/v9"
	"strings"
)

type Storage struct {
	client     redis.UniversalClient
	codec      codec.Codec[task.Task, []byte]
	hashMapKey string
}

func New(client redis.UniversalClient, codec codec.Codec[task.Task, []byte], hashMapKey string) *Storage {
	return &Storage{
		client:     client,
		codec:      codec,
		hashMapKey: hashMapKey,
	}
}

func (storage *Storage) Create(ctx context.Context, _task task.Task) (task.Task, error) {
	_task.ID = storage.generateID()
	encodedResource, err := storage.codec.Encode(_task)
	if err != nil {
		return task.Task{}, fmt.Errorf("failed to encode: %w", err)
	}

	result, err := storage.client.HSetNX(ctx, storage.hashMapKey, _task.ID, encodedResource).Result()
	if err != nil {
		return task.Task{}, fmt.Errorf("failed to execute HSetNX: %w", err)
	}

	if !result {
		return task.Task{}, errors.New("already initialized")
	}

	return _task, nil
}

func (storage *Storage) Delete(ctx context.Context, ids ...string) error {
	_, err := storage.client.HDel(ctx, storage.hashMapKey, ids...).Result()
	if err != nil {
		return fmt.Errorf("failed to execute HDel: %w", err)
	}

	return nil
}

func (storage *Storage) Load(ctx context.Context, ids ...string) (result []task.Task, err error) {
	if len(ids) == 0 {
		return nil, nil
	}

	values, err := storage.client.HMGet(ctx, storage.hashMapKey, ids...).Result()
	if err != nil {
		return nil, fmt.Errorf("failed to execute hmget: %w", err)
	}

	for idx, value := range values {
		id := ids[idx]
		_task, err := storage.parseTaskFromRedisValue(id, value)
		if err != nil {
			return nil, fmt.Errorf("failed to parse a task (id='%s'): %w", id, err)
		}

		result = append(result, _task)
	}

	return result, nil
}

func (storage *Storage) parseTaskFromRedisValue(id string, value any) (_task task.Task, err error) {
	if value == nil {
		_task = task.Task{
			Inputs:     nil,
			Outputs:    nil,
			Properties: task.Properties{},
			ID:         id,
		}
	} else {
		encodedResource, ok := value.(string)
		if !ok {
			return _task, errors.New("redis value is not a string")
		}

		_task, err = storage.codec.Decode([]byte(encodedResource))
		if err != nil {
			return
		}
	}

	return
}

func (storage *Storage) generateID() string {
	return strings.ReplaceAll(uuid.New().String(), "-", "")
}

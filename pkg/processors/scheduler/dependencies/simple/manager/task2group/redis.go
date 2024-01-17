package task2group

import (
	"context"
	"github.com/redis/go-redis/v9"
)

type RedisStorage struct {
	Client redis.UniversalClient
}

func (storage *RedisStorage) Save(ctx context.Context, task string, group string) error {
	// Save task to group mapping
	if err := storage.Client.HSet(ctx, "task2group", task, group).Err(); err != nil {
		return err
	}

	// Save group to task mapping
	if err := storage.Client.HSet(ctx, "group2task", group, task).Err(); err != nil {
		return err
	}

	return nil
}

func (storage *RedisStorage) TaskByGroup(ctx context.Context, group string) (task string, err error) {
	return storage.Client.HGet(ctx, "group2task", group).Result()
}

func (storage *RedisStorage) GroupByTask(ctx context.Context, task string) (group string, err error) {
	return storage.Client.HGet(ctx, "task2group", task).Result()
}

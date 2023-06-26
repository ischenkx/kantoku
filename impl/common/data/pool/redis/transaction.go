package redipool

import (
	"context"
	"go/types"
	"kantoku/common/data/transactional"
)

var _ transactional.Object[types.Object] = &Transaction[types.Object]{}

type Transaction[T any] struct {
	data T
	pool *Pool[T]
}

func (t *Transaction[T]) Get(ctx context.Context) (T, error) {
	return t.data, nil
}

func (t *Transaction[T]) Commit(ctx context.Context) error {
	return nil
}

func (t *Transaction[T]) Rollback(ctx context.Context) error {
	return t.pool.client.LPush(ctx, t.pool.topicName, t.data).Err()
}

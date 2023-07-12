package redipool

import (
	"context"
	"errors"
	"go/types"
	"kantoku/common/data/transactional"
)

var _ transactional.Object[types.Object] = &Transaction[types.Object]{}

type Transaction[T any] struct {
	data     T
	finished bool
	pool     *Pool[T]
}

func (t *Transaction[T]) Get(_ context.Context) (T, error) {
	return t.data, nil
}

func (t *Transaction[T]) Commit(_ context.Context) error {
	t.finished = true
	return nil
}

func (t *Transaction[T]) Rollback(ctx context.Context) error {
	if t.finished {
		return errors.New("transaction has been finished")
	}
	t.finished = true
	return t.pool.client.LPush(ctx, t.pool.topicName, t.data).Err()
}

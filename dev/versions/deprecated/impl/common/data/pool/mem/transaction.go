package mempool

import (
	"context"
	"go/types"
	"kantoku/common/data/transactional"
)

var _ transactional.Object[types.Object] = &Transaction[types.Object]{}

type Transaction[T any] struct {
	data    T
	success chan bool
}

func (t *Transaction[T]) Get(ctx context.Context) (T, error) {
	return t.data, nil
}

func (t *Transaction[T]) Commit(ctx context.Context) error {
	t.success <- true
	return nil
}

func (t *Transaction[T]) Rollback(ctx context.Context) error {
	t.success <- false
	return nil
}

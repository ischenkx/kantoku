package kvpool

import (
	"context"
	"kantoku/common/data"
	"kantoku/common/data/kv"
)

type Writer[T data.Identifiable] struct {
	storage kv.Writer[T]
}

func NewWriter[T data.Identifiable](db kv.Writer[T]) *Writer[T] {
	return &Writer[T]{
		storage: db,
	}
}

func (w *Writer[T]) Write(ctx context.Context, item T) error {
	id := item.ID()
	_, err := w.storage.Set(ctx, id, item)
	return err
}

package cellpool

import (
	"context"
	"kantoku/framework/cell"
)

// TODO: T must be data.Identifiable
type Writer[T any] struct {
	storage    cell.Storage[T]
	identifier func(context.Context, T) (string, error)
}

func NewWriter[T any](
	storage cell.Storage[T],
	identifier func(context.Context, T) (string, error)) *Writer[T] {
	return &Writer[T]{
		storage:    storage,
		identifier: identifier,
	}
}

func (w *Writer[T]) Write(ctx context.Context, item T) error {
	id, err := w.identifier(ctx, item)
	if err != nil {
		return err
	}
	return w.storage.Set(ctx, id, item)
}

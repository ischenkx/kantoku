package funcpool

import (
	"context"
)

type Writer[T any] struct {
	f func(ctx context.Context, item T) error
}

func NewWriter[T any](f func(ctx context.Context, item T) error) Writer[T] {
	return Writer[T]{f: f}
}

func (w Writer[T]) Write(ctx context.Context, item T) error {
	return w.f(ctx, item)
}

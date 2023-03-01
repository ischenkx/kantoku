package queue

import "context"

type Queue[T any] interface {
	Put(ctx context.Context, item T) error
	Clear(ctx context.Context) error
	Read(ctx context.Context) (<-chan T, error)
}

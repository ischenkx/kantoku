package db

import "context"

type KV[T any] interface {
	Get(ctx context.Context, id string) (T, error)
	Delete(ctx context.Context, id string) error
	Create(ctx context.Context, id string, item T) (T, error)
}

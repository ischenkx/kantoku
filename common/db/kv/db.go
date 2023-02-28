package kv

import "context"

type Database[T any] interface {
	Set(ctx context.Context, id string, item T) (T, error)
	Get(ctx context.Context, id string) (T, error)
	Del(ctx context.Context, id string) error
}

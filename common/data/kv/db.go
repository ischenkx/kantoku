package kv

import "context"

type Getter[T any] interface {
	Get(ctx context.Context, id string) (T, error)
}

type Setter[T any] interface {
	Set(ctx context.Context, id string, item T) (T, error)
}

type Deleter interface {
	Del(ctx context.Context, id string) error
}

type Reader[T any] interface {
	Getter[T]
}

type Writer[T any] interface {
	Setter[T]
	Deleter
}

type Database[T any] interface {
	Reader[T]
	Writer[T]
}

package kv

import "context"

type Getter[K, V any] interface {
	Get(ctx context.Context, id K) (V, error)
}

type Setter[K, V any] interface {
	Set(ctx context.Context, id K, item V) error
	GetOrSet(ctx context.Context, id K, item V) (value V, set bool, err error)
}

type Deleter[K any] interface {
	Del(ctx context.Context, id K) error
}

type Reader[K, V any] interface {
	Getter[K, V]
}

type Writer[K, V any] interface {
	Setter[K, V]
	Deleter[K]
}

type Database[K, V any] interface {
	Reader[K, V]
	Writer[K, V]
}

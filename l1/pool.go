package l1

import "context"

type PoolReader[Item any] interface {
	Channel(ctx context.Context) <-chan Item
}

type PoolWriter[Item any] interface {
	Put(ctx context.Context, item Item) error
}

type Pool[Item any] interface {
	PoolReader[Item]
	PoolWriter[Item]
}

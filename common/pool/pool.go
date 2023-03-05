package pool

import "context"

type Reader[Item any] interface {
	Read(ctx context.Context) (<-chan Item, error)
}

type Writer[Item any] interface {
	Write(ctx context.Context, item Item) error
}

type Pool[Item any] interface {
	Reader[Item]
	Writer[Item]
}

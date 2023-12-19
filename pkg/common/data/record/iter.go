package record

import "context"

type Iter[Item any] interface {
	Next(ctx context.Context) (Item, error)
	Close(ctx context.Context) error
}

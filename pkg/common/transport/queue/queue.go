package queue

import "context"

type Queue[Item any] interface {
	Consumer[Item]
	Publisher[Item]
}

type Consumer[Item any] interface {
	Consume(ctx context.Context, group string) (<-chan Message[Item], error)
}

type Publisher[Item any] interface {
	Publish(ctx context.Context, item Item) error
}

type FunctionalPublisher[Item any] struct {
	Func func(ctx context.Context, item Item) error
}

func (f FunctionalPublisher[Item]) Publish(ctx context.Context, item Item) error {
	if f.Func == nil {
		return nil
	}

	return f.Func(ctx, item)
}

package event

import (
	"context"
	"io"
)

type Listener interface {
	Listen(ctx context.Context, topics ...string) (<-chan Event, error)
}

type Publisher interface {
	Publish(ctx context.Context, events ...Event) error
}

type Bus interface {
	Listener
	Publisher
	io.Closer
}

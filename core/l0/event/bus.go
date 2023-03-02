package event

import (
	"context"
)

type Listener interface {
	Listen(ctx context.Context, topics ...string) (<-chan Event, error)
}

type Publisher interface {
	Publish(ctx context.Context, event Event) error
}

// Bus is supposed to be a fan-out queue that is used
// for event processing
type Bus interface {
	Listener
	Publisher
}

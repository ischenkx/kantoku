package platform

import "context"

type Event struct {
	Name  string
	Topic string
	Data  []byte
}

// Broker is supposed to be a fan-out queue that is used
// for event processing
type Broker interface {
	Listen() Listener
	Publish(ctx context.Context, event Event) error
}

type Listener interface {
	Subscribe(ctx context.Context, topics ...string) error
	Unsubscribe(ctx context.Context, topics ...string) error
	UnsubscribeAll(ctx context.Context) error
	Incoming(ctx context.Context) (<-chan Event, error)
	Close(ctx context.Context) error
}

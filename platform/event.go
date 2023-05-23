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
	Listen(ctx context.Context, topics ...string) (<-chan Event, error)
	Publish(ctx context.Context, event Event) error
}

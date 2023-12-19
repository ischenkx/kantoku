package event

import "context"

type Queue struct {
	Name   string
	Topics []string
}

type Bus interface {
	Consume(ctx context.Context, queue Queue) (<-chan Event, error)
	Publish(ctx context.Context, events Event) error
}

package client

import (
	"context"
	"errors"
	event2 "github.com/ischenkx/kantoku/pkg/system/kernel/event"
)

type eventBus struct{}

func (e eventBus) Consume(ctx context.Context, queue event2.Queue) (<-chan event2.Event, error) {
	return nil, errors.New("not supported by an http client")
}

func (e eventBus) Publish(ctx context.Context, events event2.Event) error {
	return errors.New("not supported by an http client")
}

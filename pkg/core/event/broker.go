package event

import (
	"context"
	"github.com/ischenkx/kantoku/pkg/common/broker"
)

type Broker struct {
	broker broker.Broker[Event]
}

func NewBroker(b broker.Broker[Event]) *Broker {
	return &Broker{broker: b}
}

func (b *Broker) Send(ctx context.Context, event Event) error {
	return b.broker.Publish(ctx, event.Topic, event)
}

func (b *Broker) Consume(ctx context.Context, info broker.TopicsInfo) (<-chan broker.Message[Event], error) {
	return b.broker.Consume(ctx, info)
}

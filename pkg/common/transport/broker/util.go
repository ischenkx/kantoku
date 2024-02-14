package broker

import (
	"context"
	"github.com/ischenkx/kantoku/pkg/common/transport/queue"
)

func Restrict[Item any](broker Broker[Item], topic string) RestrictedBroker[Item] {
	return RestrictedBroker[Item]{
		Topic:  topic,
		Broker: broker,
	}
}

type RestrictedBroker[Item any] struct {
	Topic  string
	Broker Broker[Item]
}

func (queue RestrictedBroker[Item]) Consume(ctx context.Context, group string) (<-chan queue.Message[Item], error) {
	return queue.Broker.Consume(ctx, TopicsInfo{
		Group:  group,
		Topics: []string{queue.Topic},
	})
}

func (queue RestrictedBroker[Item]) Publish(ctx context.Context, item Item) error {
	return queue.Broker.Publish(ctx, queue.Topic, item)
}

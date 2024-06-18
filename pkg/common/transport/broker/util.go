package broker

import (
	"context"
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

func (queue RestrictedBroker[Item]) Consume(ctx context.Context, group string) (<-chan Message[Item], error) {
	return queue.Broker.Consume(ctx, TopicsInfo{
		Group:  group,
		Topics: []string{queue.Topic},
	})
}

func (queue RestrictedBroker[Item]) Publish(ctx context.Context, item Item) error {
	return queue.Broker.Publish(ctx, queue.Topic, item)
}

type HandlerFunc[Item any] func(ctx context.Context, ev Item) error

func Process[Item any](ctx context.Context, message Message[Item], handler HandlerFunc[Item]) error {
	defer message.Nack()
	if handler != nil {
		if err := handler(ctx, message.Item()); err != nil {
			return err
		}
	}
	message.Ack()
	return nil
}

type Processor[Item any] struct {
	Handler      HandlerFunc[Item]
	ErrorHandler func(ctx context.Context, ev Item, err error)
}

func (processor Processor[Item]) Process(ctx context.Context, channel <-chan Message[Item]) {
	for {
		select {
		case <-ctx.Done():
			return
		case message := <-channel:
			if err := Process[Item](ctx, message, processor.Handler); err != nil {
				if processor.ErrorHandler != nil {
					processor.ErrorHandler(ctx, message.Item(), err)
				}
			}
		}
	}
}

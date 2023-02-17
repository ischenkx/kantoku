package redis

import (
	"context"
	"github.com/redis/go-redis/v9"
	"hayaku/common/codec"
	"hayaku/l0/event"
	"log"
	"strings"
)

type Bus struct {
	client redis.UniversalClient
	codec  codec.Codec[event.Event]
}

func NewBus(client redis.UniversalClient, codec codec.Codec[event.Event]) *Bus {
	return &Bus{
		client: client,
		codec:  codec,
	}
}

func (b *Bus) Listen(ctx context.Context, topics ...string) (<-chan event.Event, error) {
	channel := make(chan event.Event)
	go func(outputs chan<- event.Event) {
		pubsub := b.client.Subscribe(ctx, topics...)
		messages := pubsub.Channel()

		defer close(outputs)
		defer pubsub.Close()

	processor:
		for {
			select {
			case message := <-messages:
				ev, err := b.codec.Decode(strings.NewReader(message.Payload))
				if err != nil {
					log.Println("failed to decode the event:", err)
					continue
				}
				outputs <- ev
			case <-ctx.Done():
				break processor
			}
		}
	}(channel)
	return channel, nil
}

func (b *Bus) Publish(ctx context.Context, event event.Event) error {
	payload, err := b.codec.Encode(event)
	if err != nil {
		return nil
	}
	return b.client.Publish(ctx, event.Topic, payload).Err()
}

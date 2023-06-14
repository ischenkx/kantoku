package redivent

import (
	"context"
	"github.com/redis/go-redis/v9"
	"kantoku/common/codec"
	"kantoku/kernel/platform"
)

type Broker struct {
	codec  codec.Codec[platform.Event, []byte]
	client redis.UniversalClient
}

func New(codec codec.Codec[platform.Event, []byte], client redis.UniversalClient) *Broker {
	return &Broker{
		codec:  codec,
		client: client,
	}
}

func (b *Broker) Listen() platform.Listener {
	return NewListener(b.codec, b.client.Subscribe(context.Background()))
}

func (b *Broker) Publish(ctx context.Context, event platform.Event) error {
	message, err := b.codec.Encode(event)
	if err != nil {
		return err
	}

	return b.client.Publish(ctx, event.Topic, message).Err()
}

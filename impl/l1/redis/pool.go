package redis

import (
	"bytes"
	"context"
	"github.com/redis/go-redis/v9"
	"kantoku/common/codec"
	"log"
)

type Pool[Item any] struct {
	codec codec.Codec[Item]
	redis redis.UniversalClient
	topic string
}

func NewPool[Item any](codec codec.Codec[Item], redisClient redis.UniversalClient, topic string) *Pool[Item] {
	return &Pool[Item]{
		codec: codec,
		redis: redisClient,
		topic: topic,
	}
}

func (p *Pool[Item]) Channel(ctx context.Context) <-chan Item {
	channel := make(chan Item, 512)
	pubsub := p.redis.Subscribe(ctx, p.topic)

	go func(ctx context.Context, pubsub *redis.PubSub, outputs chan<- Item) {
		redisChannel := pubsub.Channel()

		defer pubsub.Close()
		defer close(outputs)

	processor:
		for {
			select {
			case <-ctx.Done():
				break processor
			case message := <-redisChannel:
				raw := []byte(message.Payload)
				item, err := p.codec.Decode(bytes.NewReader(raw))
				if err != nil {
					log.Printf("failed to decode the incoming message: %s\n", err)
					continue
				}
				outputs <- item
			}
		}
	}(ctx, pubsub, channel)

	return channel
}

func (p *Pool[Item]) Put(ctx context.Context, item Item) error {
	raw, err := p.codec.Encode(item)
	if err != nil {
		return err
	}

	cmd := p.redis.Publish(ctx, p.topic, raw)
	return cmd.Err()
}

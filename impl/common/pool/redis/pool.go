package redipool

import (
	"context"
	"github.com/redis/go-redis/v9"
	"kantoku/common/chutil"
	"kantoku/common/codec"
	"log"
	"strings"
)

type Pool[T any] struct {
	client    redis.UniversalClient
	codec     codec.Codec[T]
	topicName string
}

func New[T any](client redis.UniversalClient, codec codec.Codec[T], topicName string) *Pool[T] {
	return &Pool[T]{
		client:    client,
		codec:     codec,
		topicName: topicName,
	}
}

func (pool *Pool[T]) Write(ctx context.Context, item T) error {
	data, err := pool.codec.Encode(item)
	if err != nil {
		return err
	}

	cmd := pool.client.Publish(ctx, pool.topicName, data)
	if cmd.Err() != nil {
		return cmd.Err()
	}
	return nil
}

func (pool *Pool[T]) Read(ctx context.Context) (<-chan T, error) {
	pubsub := pool.client.Subscribe(ctx, pool.topicName)
	channel := make(chan T, 1024)
	chutil.SyncWithContext(ctx, channel)

	go func(ctx context.Context, ps *redis.PubSub, outputs chan<- T) {
		pubsubChannel := ps.Channel()
		defer ps.Close()

	loop:
		for {
			select {
			case message := <-pubsubChannel:
				data, err := pool.codec.Decode(strings.NewReader(message.Payload))
				if err != nil {
					log.Println("failed to decode the incoming message:", err)
				}
				outputs <- data
			case <-ctx.Done():
				break loop
			}
		}
	}(ctx, pubsub, channel)

	return channel, nil
}

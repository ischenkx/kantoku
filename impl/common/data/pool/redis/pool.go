package redipool

import (
	"context"
	"github.com/redis/go-redis/v9"
	"kantoku/common/chutil"
	"kantoku/common/codec"
	"log"
)

type Pool[T any] struct {
	client    redis.UniversalClient
	codec     codec.Codec[T, []byte]
	topicName string
}

func New[T any](client redis.UniversalClient, codec codec.Codec[T, []byte], topicName string) *Pool[T] {
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

	cmd := pool.client.RPush(ctx, pool.topicName, data)
	if cmd.Err() != nil {
		return cmd.Err()
	}
	return nil
}

func (pool *Pool[T]) Read(ctx context.Context) (<-chan T, error) {
	channel := make(chan T, 1024)
	chutil.CloseWithContext(ctx, channel)

	go func(ctx context.Context, outputs chan<- T) {
	loop:
		for {
			select {
			case <-ctx.Done():
				break loop
			default:
			}
			result, err := pool.client.BLPop(ctx, 0, pool.topicName).Result()
			if err != nil {
				log.Println("failed to pop a task from the queue:", err)
				continue
			}

			if len(result) != 2 {
				log.Println("length of the result is not equal 2")
				continue
			}

			data := result[1]

			output, err := pool.codec.Decode([]byte(data))
			if err != nil {
				log.Println("failed to decode the incoming message:", err)
				continue
			}
			outputs <- output
		}
	}(ctx, channel)

	return channel, nil
}

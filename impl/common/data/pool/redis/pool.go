package redipool

import (
	"context"
	"fmt"
	"github.com/redis/go-redis/v9"
	"github.com/samber/lo"
	"go/types"
	"kantoku/common/chutil"
	"kantoku/common/codec"
	"kantoku/common/data/pool"
	"kantoku/common/data/transactional"
	"log"
	"time"
)

var _ pool.Pool[types.Object] = &Pool[types.Object]{}

// Pool with redis. This implementation does not guarantee FIFO because it never blocks queue.
type Pool[T any] struct {
	client    redis.UniversalClient
	codec     codec.Codec[T, []byte]
	topicName string
}

const PollPeriod = time.Millisecond * 100

func New[T any](client redis.UniversalClient, codec codec.Codec[T, []byte], topicName string) *Pool[T] {
	return &Pool[T]{
		client:    client,
		codec:     codec,
		topicName: topicName,
	}
}

func (pool *Pool[T]) Write(ctx context.Context, items ...T) error {
	var err error
	data := make([][]byte, len(items))
	for i, item := range items {
		data[i], err = pool.codec.Encode(item)
		if err != nil {
			break
		}
	}
	if err != nil {
		return fmt.Errorf("failed to encode items: %w", err)
	}

	cmd := pool.client.RPush(ctx, pool.topicName,
		lo.Map(data, func(item []byte, _ int) interface{} { return item })...)

	if cmd.Err() != nil {
		return cmd.Err()
	}
	return nil
}

func (pool *Pool[T]) Read(ctx context.Context) (<-chan transactional.Object[T], error) {
	channel := make(chan transactional.Object[T], 0)
	chutil.CloseWithContext(ctx, channel)

	go func(ctx context.Context, outputs chan<- transactional.Object[T]) {
	loop:
		for {
			select {
			case <-ctx.Done():
				break loop
			default:
			}

			result, err := pool.client.LPop(ctx, pool.topicName).Result()

			if err == redis.Nil {
				// nothing in queue
				time.Sleep(PollPeriod)
				continue
			}

			if err != nil {
				log.Println("failed to pop a task from the queue:", err)
				continue
			}
			data := []byte(result)
			output, err := pool.codec.Decode(data)
			if err != nil {
				log.Println("failed to decode the incoming message:", err)
				continue
			}

			select {
			case <-ctx.Done():
				// push results back value if context is closed
				if err := pool.client.LPush(context.Background(), pool.topicName, data).Err(); err != nil {
					log.Println("failed to push back value which was read after context was closed:", err)
				}
				break loop
			case outputs <- &Transaction[T]{
				data: output,
				pool: pool,
			}:
			}

		}
	}(ctx, channel)

	return channel, nil
}

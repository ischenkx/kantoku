package queue

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"github.com/redis/go-redis/v9"
	"kantoku/common/codec"
	"log"
)

// RedisQueue is a queue implementation that uses Redis as a backend.
type RedisQueue[T any] struct {
	client redis.UniversalClient
	key    string
	codec  codec.Codec[T]
}

// NewRedisQueue returns a new RedisQueue instance.
func NewRedisQueue[T any](client redis.UniversalClient, key string, codec codec.Codec[T]) *RedisQueue[T] {
	return &RedisQueue[T]{client: client, key: key, codec: codec}
}

// Put adds an item to the queue.
func (q *RedisQueue[T]) Put(ctx context.Context, item T) error {
	encoded, err := q.codec.Encode(item)
	if err != nil {
		return err
	}

	_, err = q.client.RPush(ctx, q.key, encoded).Result()
	return err
}

func (q *RedisQueue[T]) Clear(ctx context.Context) error {
	_, err := q.client.Del(ctx, q.key).Result()
	return err
}

// Read returns a channel that streams items from the queue.
func (q *RedisQueue[T]) Read(ctx context.Context) (<-chan T, error) {
	stream := make(chan T)
	go func() {
		defer close(stream)
		for {
			select {
			case <-ctx.Done():
				return
			default:
				result, err := q.client.BLPop(ctx, 0, q.key).Result()
				if err != nil {
					if !errors.Is(err, redis.Nil) {
						// Log the error if it's not just a timeout.
						log.Println(fmt.Errorf("BLPop error: %v", err))
					}
					continue
				}

				if len(result) != 2 {
					log.Println(errors.New("invalid result"))
				}

				data := result[1]
				item, err := q.codec.Decode(bytes.NewReader([]byte(data)))
				if err != nil {
					log.Println(fmt.Errorf("decoder error: %v", err))
				}

				select {
				case stream <- item:
				case <-ctx.Done():
					return
				}
			}
		}
	}()
	return stream, nil
}

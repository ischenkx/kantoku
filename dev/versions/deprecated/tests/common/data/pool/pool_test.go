package pool

import (
	"context"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
	"kantoku/common/data/pool"
	"kantoku/impl/common/codec/jsoncodec"
	mempool "kantoku/impl/common/data/pool/mem"
	redipool "kantoku/impl/common/data/pool/redis"
	"testing"
	"time"
)

func newRedisPool[Item any](ctx context.Context) pool.Pool[Item] {
	client := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379", // Redis server address
		Password: "",               // Redis server password (leave empty if not set)
		DB:       0,                // Redis database index
	})

	if cmd := client.Ping(ctx); cmd.Err() != nil {
		panic("failed to ping the redis client: " + cmd.Err().Error())
	}
	if cmd := client.Del(ctx, "TEST_POOL"); cmd.Err() != nil {
		panic("failed to clear topic: " + cmd.Err().Error())
	}
	return redipool.New[Item](client, jsoncodec.New[Item](), "TEST_POOL")
}

func newMemPool[Item any](_ context.Context) pool.Pool[Item] {
	return mempool.New[Item](mempool.DefaultConfig)
}

func TestPool(t *testing.T) {
	implementations := map[string]func(context.Context) pool.Pool[string]{
		"mem":   newMemPool[string],
		"redis": newRedisPool[string],
	}

	for label, impl := range implementations {
		t.Run(label+": PutNothingAndGetNothing", func(t *testing.T) {
			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()

			p := impl(ctx)

			itemsCh, err := p.Read(ctx)
			assert.NoError(t, err)

			select {
			case <-itemsCh:
				t.Error("Expected no items, but received an item")
			case <-time.After(3 * time.Second):
				// Passed, no items received within the timeout
			}
		})

		t.Run(label+": PutOneItemGetItCommit", func(t *testing.T) {
			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()
			p := impl(ctx)

			err := p.Write(ctx, "item1")
			assert.NoError(t, err)

			itemsCh, err := p.Read(ctx)
			assert.NoError(t, err)

			select {
			case tx := <-itemsCh:
				item, err := tx.Get(ctx)
				assert.NoError(t, err)
				assert.Equal(t, "item1", item)
				err = tx.Commit(ctx)
				assert.NoError(t, err)
			case <-time.After(3 * time.Second):
				t.Error("Expected an item, but none received")
			}
		})

		t.Run(label+": PutTwoItemsGetOneRollbackGetOneCommit", func(t *testing.T) {
			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()

			p := impl(ctx)

			err := p.Write(ctx, "item1", "item2")
			assert.NoError(t, err)

			ctx1, cancel1 := context.WithCancel(context.Background())
			defer cancel1()
			itemsCh, err := p.Read(ctx1)
			assert.NoError(t, err)

			select {
			case tx := <-itemsCh:
				item, err := tx.Get(ctx1)
				assert.NoError(t, err)
				assert.Equal(t, "item1", item)
				err = tx.Rollback(ctx1)
				assert.NoError(t, err)
			case <-time.After(3 * time.Second):
				t.Error("Expected an item, but none received")
			}
			cancel1()

			ctx2, cancel2 := context.WithCancel(context.Background())
			defer cancel2()
			itemsCh, err = p.Read(ctx2)
			assert.NoError(t, err)

			select {
			case tx := <-itemsCh:
				item, err := tx.Get(ctx2)
				assert.NoError(t, err)
				assert.True(t, item == "item1" || item == "item2") // order is not guaranteed anymore
				err = tx.Commit(ctx2)
				assert.NoError(t, err)
			case <-time.After(3 * time.Second):
				t.Error("Expected an item, but none received")
			}
			cancel2()
		})

		t.Run(label+": PutRandomNumbersGetAndCommit", func(t *testing.T) {
			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()

			p := impl(ctx)
			itemsCh, err := p.Read(ctx)
			assert.NoError(t, err)

			for i := 0; i <= 10; i++ {
				// Generate a random number
				generated := uuid.New().String()

				err := p.Write(ctx, generated)
				assert.NoError(t, err)

				select {
				case tx := <-itemsCh:
					received, err := tx.Get(ctx)
					assert.NoError(t, err)
					assert.Equal(t, generated, received)
					err = tx.Commit(ctx)
					assert.NoError(t, err)
				case <-time.After(3 * time.Second):
					t.Error("Expected an item, but none received")
				}
			}
		})
	}
}

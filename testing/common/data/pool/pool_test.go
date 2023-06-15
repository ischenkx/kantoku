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

type Item struct {
	Data string
	Name string
}

func newRedisPool(ctx context.Context) *redipool.Pool[Item] {
	client := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379", // Redis server address
		Password: "",               // Redis server password (leave empty if not set)
		DB:       0,                // Redis database index
	})

	if cmd := client.Ping(ctx); cmd.Err() != nil {
		panic("failed to ping the redis client: " + cmd.Err().Error())
	}

	return redipool.New[Item](client, jsoncodec.New[Item](), "TEST_POOL")
}

func newMemPool[Item any](ctx context.Context) pool.Pool[Item] {
	return mempool.New[Item](mempool.DefaultConfig)
}

func TestPool(t *testing.T) {
	ctx := context.Background()
	implementations := map[string]func(context.Context) pool.Pool[string]{
		"mem": newMemPool[string],
	}

	for label, impl := range implementations {
		t.Run(label+": PutNothingAndGetNothing", func(t *testing.T) {
			p := impl(ctx)

			ctx := context.Background()
			itemsCh, err := p.Read(ctx)
			assert.NoError(t, err)

			select {
			case <-itemsCh:
				t.Error("Expected no items, but received an item")
			case <-time.After(3 * time.Second):
				// Passed, no items received within the timeout
			}
			t.Log("finished")
		})

		t.Run(label+": PutOneItemGetItCommit", func(t *testing.T) {
			p := impl(ctx)

			ctx := context.Background()
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
			t.Log("finished")
		})

		t.Run(label+": PutTwoItemsGetOneRollbackGetOneCommit", func(t *testing.T) {
			p := impl(ctx)

			ctx := context.Background()
			err := p.Write(ctx, "item1", "item2")
			assert.NoError(t, err)

			itemsCh, err := p.Read(ctx)
			assert.NoError(t, err)

			select {
			case tx := <-itemsCh:
				item, err := tx.Get(ctx)
				assert.NoError(t, err)
				assert.Equal(t, "item1", item)
				err = tx.Rollback(ctx)
				assert.NoError(t, err)
			case <-time.After(3 * time.Second):
				t.Error("Expected an item, but none received")
			}

			itemsCh, err = p.Read(ctx)
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
			t.Log("finished")
		})

		t.Run(label+": PutRandomNumbersGetAndCommit", func(t *testing.T) {
			p := impl(ctx)

			ctx := context.Background()

			for i := 1; i <= 10; i++ {
				// Generate a random number
				generated := uuid.New().String()

				err := p.Write(ctx, generated)
				assert.NoError(t, err)

				itemsCh, err := p.Read(ctx)
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
			t.Log("finished")
		})
	}
}

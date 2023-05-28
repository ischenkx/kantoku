package pool

import (
	"context"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	"github.com/samber/lo"
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

func newMemPool(ctx context.Context) *mempool.Pool[Item] {
	return mempool.New[Item]()
}

func TestPool(t *testing.T) {
	ctx := context.Background()
	implementations := map[string]pool.Pool[Item]{
		"redis":  newRedisPool(ctx),
		"memory": newMemPool(ctx),
	}

	for label, impl := range implementations {
		t.Run(label, func(t *testing.T) {
			var testSet []Item

			for i := 0; i < 10; i++ {
				testSet = append(testSet, Item{
					Data: uuid.New().String(),
					Name: uuid.New().String(),
				})
			}

			lo.ForEach(testSet, func(item Item, _ int) {
				if err := impl.Write(ctx, item); err != nil {
					t.Fatal("failed to write in pool:", err)
				}
			})

			cancelableContext, cancel := context.WithCancel(ctx)

			channel, err := impl.Read(cancelableContext)
			if err != nil {
				t.Fatal("failed to read from a pool:", err)
			}
			time.Sleep(time.Second * 3)
			cancel()

			receivedItems := lo.ChannelToSlice[Item](channel)

			assert.Equal(t, testSet, receivedItems, "test_set must be equal to a received set from the pool")
		})
	}
}

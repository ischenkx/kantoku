package kv

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	"github.com/samber/lo"
	"github.com/stretchr/testify/assert"
	"kantoku/common/data"
	"kantoku/common/data/kv"
	"kantoku/impl/common/codec/jsoncodec"
	redikv "kantoku/impl/common/data/kv/redis"
	"testing"
)

type Item struct {
	Data string
	Name string
}

func newRedisKV(ctx context.Context) *redikv.DB[Item] {
	client := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379", // Redis server address
		Password: "",               // Redis server password (leave empty if not set)
		DB:       0,                // Redis database index
	})

	if cmd := client.Ping(ctx); cmd.Err() != nil {
		panic("failed to ping the redis client: " + cmd.Err().Error())
	}

	return redikv.New[Item](client, jsoncodec.New[Item](), "TEST_KV")
}

func TestKV(t *testing.T) {
	ctx := context.Background()
	implementations := map[string]kv.Database[string, Item]{
		"redis": newRedisKV(ctx),
	}

	for label, impl := range implementations {
		t.Run(fmt.Sprintf("%s Get/Set/Del", label), func(t *testing.T) {
			testSet := map[string]Item{}

			for i := 0; i < 10; i++ {
				testSet[uuid.New().String()] = Item{Data: uuid.New().String(), Name: uuid.New().String()}
			}

			for key, item := range testSet {
				if err := impl.Set(ctx, key, item); err != nil {
					t.Fatal("failed to set:", err)
				}
			}

			for key, item := range testSet {
				received, err := impl.Get(ctx, key)
				if err != nil {
					t.Fatalf("failed to get an item (key='%s'): %s", key, err)

				}

				assert.Equal(t, item, received)

				if err := impl.Del(ctx, key); err != nil {
					t.Fatalf("failed to delete an item by a key (key='%s'): %s", key, err)
				}

				testSet[key] = Item{}

				receivedSet := lo.SliceToMap(
					lo.Keys(testSet),
					func(key string) (string, Item) {
						item, err := impl.Get(ctx, key)
						if err != nil {
							if err == data.NotFoundErr {
								item = Item{}
							} else {
								t.Fatal("failed to get an item:", err)
							}
						}

						return key, item
					},
				)

				assert.Equal(t, testSet, receivedSet, "test set must be equal to a received set")
			}

		})

		t.Run(fmt.Sprintf("%s GetOrSet", label), func(t *testing.T) {
			key := uuid.New().String()
			val1 := Item{Data: uuid.New().String(), Name: uuid.New().String()}
			val2 := Item{Data: uuid.New().String(), Name: uuid.New().String()}

			receivedValue1, _, err := impl.GetOrSet(ctx, key, val1)
			if err != nil {
				t.Fatal("failed:", err)
			}
			assert.Equal(t, val1, receivedValue1)

			receivedValue2, _, err := impl.GetOrSet(ctx, key, val2)
			if err != nil {
				t.Fatal("failed:", err)
			}
			assert.Equal(t, val1, receivedValue2)
		})
	}
}

package testing

import (
	"context"
	"fmt"
	"github.com/stretchr/testify/require"
	jsonCodec "kantoku/impl/common/codec/jsoncodec"
	redisQueue "kantoku/impl/common/data/pool/redis"
	"kantoku/testing/common"
	"testing"
)

type testStruct struct {
	Name   string `json:"name"`
	Number int    `json:"number"`
}

func TestRedisQueue(t *testing.T) {
	// Set up Redis client and codec.
	client := common.DefaultClient()
	codec := jsonCodec.Codec[testStruct]{}

	t.Run("publish-read", func(t *testing.T) {
		q := redisQueue.New[testStruct](client, codec, t.Name())

		// Subscribe to the pool
		ch, err := q.Read(context.Background())
		require.NoError(t, err)

		// Add an item to the pool.
		item := testStruct{"test-item", 42}
		err = q.Write(context.Background(), item)
		require.NoError(t, err, "error on put")

		// Decode the item and check its value.
		decoded, ok := <-ch
		require.True(t, ok)
		require.Equal(t, item, decoded)
	})

	t.Run("add-read-many", func(t *testing.T) {
		q := redisQueue.New[testStruct](client, codec, t.Name())

		// Subscribe to the pool
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()
		stream, err := q.Read(ctx)
		require.NoError(t, err, "error on read")

		const itemsCount = 100

		// Add items to the pool.
		items := make([]testStruct, itemsCount)
		for i := 0; i < itemsCount; i++ {
			items[i] = testStruct{
				Name:   fmt.Sprintf("test-item-%d", i),
				Number: i,
			}
			err := q.Write(context.Background(), items[i])
			require.NoError(t, err, "error on add")
		}

		for i := 0; i < itemsCount; i++ {
			item, ok := <-stream
			require.True(t, ok, "Stream closed prematurely")
			require.Equal(t, items[i], item, "Queue returned wrong item")
		}
	})

	//t.Run("clear", func(t *testing.T) {
	//	client.Del(context.Background(), t.Name())
	//
	//	q := redisQueue.New[testStruct](client, codec, t.Name())
	//
	//	// Add an item to the queue.
	//	item := testStruct{"test-item", 42}
	//	err := q.Write(context.Background(), item)
	//	require.NoError(t, err)
	//	err = q.Write(context.Background(), item)
	//	require.NoError(t, err)
	//
	//	err = q.Clear(context.Background())
	//	require.NoError(t, err)
	//	// Check that queue is empty
	//	result, err := client.LRange(context.Background(), t.Name(), 0, -1).Result()
	//	require.NoError(t, err, "error on getting length")
	//	require.Len(t, result, 0)
	//})

}

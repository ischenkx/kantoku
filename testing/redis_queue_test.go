package testing

import (
	"context"
	"fmt"
	"github.com/stretchr/testify/require"
	jsonCodec "kantoku/impl/common/codec/json"
	redisQueue "kantoku/impl/common/queue"
	"kantoku/testing/common"
	"strings"
	"testing"
)

type testStruct struct {
	name   string
	number int
}

func TestRedisQueue(t *testing.T) {
	// Set up Redis client and codec.
	client := common.DefaultClient()
	codec := jsonCodec.Codec[testStruct]{}

	// Set up queue key.
	queueKey := "test-queue"

	t.Run("add", func(t *testing.T) {
		client.Del(context.Background(), queueKey)
		q := redisQueue.NewRedisQueue[testStruct](client, queueKey, codec)

		// Add an item to the queue.
		item := testStruct{"test-item", 42}
		err := q.Put(context.Background(), item)
		require.NoError(t, err, "error on put")

		// Check that the item was added to the queue.
		result, err := client.LRange(context.Background(), queueKey, 0, -1).Result()
		require.NoError(t, err, "error on getting length")
		require.Len(t, result, 1)

		// Decode the item and check its value.
		decoded, err := codec.Decode(strings.NewReader(result[0]))
		require.NoError(t, err)
		require.Equal(t, item, decoded)
	})

	t.Run("addMultiple", func(t *testing.T) {
		client.Del(context.Background(), queueKey)
		q := redisQueue.NewRedisQueue[testStruct](client, queueKey, codec)

		const itemsCount = 100

		// Add items to the queue.
		items := make([]testStruct, itemsCount)
		for i := 0; i < itemsCount; i++ {
			items[i] = testStruct{
				name:   fmt.Sprintf("test-item-%d", i),
				number: i,
			}
			err := q.Put(context.Background(), items[i])
			require.NoError(t, err, "error on add")
		}

		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()
		stream, err := q.Read(ctx)
		require.NoError(t, err, "error on read")

		// Check that the item was added to the queue.
		result, err := client.LRange(context.Background(), queueKey, 0, -1).Result()
		require.NoError(t, err)
		require.Len(t, result, itemsCount)

		for i := 0; i < itemsCount; i++ {
			item, ok := <-stream
			require.True(t, ok, "Stream closed prematurely")
			require.Equal(t, items[itemsCount-i-1], item)
		}
		// Verify that the stream is closed after all items are read.
		_, ok := <-stream
		require.False(t, ok, "Stream should be closed, but is still open")
	})
}

package resource

import (
	"context"
	"github.com/google/uuid"
	"github.com/ischenkx/kantoku/pkg/common/data/codec"
	redisResources "github.com/ischenkx/kantoku/pkg/impl/kernel/resource/redis"
	resource2 "github.com/ischenkx/kantoku/pkg/system/kernel/resource"
	"github.com/redis/go-redis/v9"
	"github.com/samber/lo"
	"github.com/stretchr/testify/assert"
	"testing"
)

type Index map[string]resource2.Resource

func newRedisResources(ctx context.Context) *redisResources.Storage {
	client := redis.NewUniversalClient(&redis.UniversalOptions{
		//Addrs: []string{"172.23.146.206:6379"},
		Addrs: []string{":6379"},
		DB:    1,
	})

	if err := client.Ping(ctx).Err(); err != nil {
		panic(err)
	}

	if _, err := client.FlushDB(ctx).Result(); err != nil {
		panic(err)
	}

	return redisResources.New(
		client,
		codec.JSON[resource2.Resource](),
		"test-resources",
	)
}

func TestResources(t *testing.T) {
	ctx := context.Background()
	implementations := map[string]resource2.Storage{
		"redis": newRedisResources(ctx),
	}

	for name, impl := range implementations {
		t.Run(name, func(t *testing.T) {
			ImplTest(ctx, t, impl)
		})
	}
}

func ImplTest(ctx context.Context, t *testing.T, impl resource2.Storage) {
	HandTest(ctx, t, impl)
}

func HandTest(ctx context.Context, t *testing.T, impl resource2.Storage) {
	index := Index{}

	for i := 0; i < 10; i++ {
		ReadTest(ctx, t, impl, index)
		for j := 0; j < 5; j++ {
			err := allocMutation(ctx, impl, index)
			assert.NoError(t, err, "failed to apply an alloc mutation")
		}
		for j := 0; j < 10; j++ {
			err := randomMutation(ctx, impl, index)
			assert.NoError(t, err, "failed to apply a random mutation")
		}
		ReadTest(ctx, t, impl, index)
	}
}

func ReadTest(ctx context.Context, t *testing.T, impl resource2.Storage, index Index) {
	checkNonExistingResources(ctx, t, impl, index)
	checkExistingResources(ctx, t, impl, index)
}

type mutation func(ctx context.Context, impl resource2.Storage, index Index) error

func randomMutation(ctx context.Context, impl resource2.Storage, index Index) error {
	return lo.Sample([]mutation{
		allocMutation,
		initMutation,
		deallocMutation,
	})(ctx, impl, index)
}

func allocMutation(ctx context.Context, impl resource2.Storage, index Index) error {
	ids, err := impl.Alloc(ctx, 100)
	if err != nil {
		return err
	}

	for _, id := range ids {
		index[id] = resource2.Resource{
			Data:   nil,
			ID:     id,
			Status: resource2.Allocated,
		}
	}

	return nil
}

func initMutation(ctx context.Context, impl resource2.Storage, index Index) error {
	ids := lo.Filter(lo.Uniq(lo.Samples(lo.Keys(index), 100)),
		func(id string, _ int) bool {
			return index[id].Status == resource2.Allocated
		})

	err := impl.Init(ctx, lo.Map(ids, func(id string, _ int) resource2.Resource {
		return resource2.Resource{
			Data: []byte("initialized"),
			ID:   id,
		}
	}))
	if err != nil {
		return err
	}

	for _, id := range ids {
		index[id] = resource2.Resource{
			Data:   []byte("initialized"),
			ID:     id,
			Status: resource2.Ready,
		}
	}

	return nil
}

func deallocMutation(ctx context.Context, impl resource2.Storage, index Index) error {
	ids := lo.Uniq(lo.Samples(lo.Keys(index), 100))

	if err := impl.Dealloc(ctx, ids); err != nil {
		return nil
	}

	for _, id := range ids {
		delete(index, id)
	}

	return nil
}

func checkExistingResources(ctx context.Context, t *testing.T, impl resource2.Storage, index Index) {
	for i := 0; i < 100; i++ {
		entries := lo.Samples(lo.Entries(index), 5)
		ids := lo.Map(entries, func(entry lo.Entry[string, resource2.Resource], _ int) string {
			return entry.Key
		})
		expectedResources := lo.Map(entries, func(entry lo.Entry[string, resource2.Resource], _ int) resource2.Resource {
			return entry.Value
		})

		loadedResources, err := impl.Load(ctx, ids...)
		assert.NoError(t, err, "failed to load a subset of ids")

		assert.Equal(t, expectedResources, loadedResources)
	}
}

func checkNonExistingResources(ctx context.Context, t *testing.T, impl resource2.Storage, index Index) {
	for i := 0; i < 100; i++ {
		id := generateId()
		if _, ok := index[id]; ok {
			continue
		}

		batch, err := impl.Load(ctx, id)
		assert.NoError(t, err, "failed to load resources")

		assert.Equal(t, resource2.DoesNotExist, batch[0].Status)
	}
}

func generateId() string {
	return uuid.New().String()
}

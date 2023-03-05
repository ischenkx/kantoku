package testing

import (
	"context"
	"fmt"
	"github.com/satori/go.uuid"
	"kantoku/impl/common/codec/jsoncodec"
	redisStorage "kantoku/impl/framework/cell/redis"
	"kantoku/testing/common"
	"testing"
)

func TestRedisStorage(t *testing.T) {
	client := common.DefaultClient()
	defer client.Close()

	t.Run("ping", func(t *testing.T) {
		name, err := client.Ping(context.Background()).Result()
		if err != nil {
			t.Fatalf("ping failed: %v\n", err)
		}
		if name != "PONG" {
			t.Fatalf("expected status PONG, actual: %v\n", name)
		}
		fmt.Printf("status: %v\n", name)
	})

	storage := redisStorage.New[string](client, jsoncodec.Codec[string]{})

	t.Run("create", func(t *testing.T) {
		id, err := storage.Make(context.Background(), "test data")
		if err != nil {
			t.Fatalf("create failed: %v", err)
		}
		if _, err := uuid.FromString(id); err != nil {
			t.Fatalf("invalid uuid: %v", err)
		}
	})

	t.Run("get", func(t *testing.T) {
		id, _ := storage.Make(context.Background(), "test data")
		c, err := storage.Get(context.Background(), id)
		if err != nil {
			t.Fatalf("get failed: %v", err)
		}
		if c.ID != id {
			t.Fatalf("incorrect id, expected %q but got %q", id, c.ID)
		}
		if string(c.Data) != "test data" {
			t.Fatalf("incorrect data, expected %q but got %q", "test data", string(c.Data))
		}
	})

	//t.Run("delete", func(t *testing.T) {
	//	id, _ := storage.Make(context.Background(), "test data")
	//	if err := storage.Delete(context.Background(), id); err != nil {
	//		t.Fatalf("delete failed: %v", err)
	//	}
	//	_, err := storage.Get(context.Background(), id)
	//	if err == nil {
	//		t.Fatalf("c should be deleted")
	//	}
	//})
}

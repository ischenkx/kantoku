package testing

import (
	"context"
	"fmt"
	"github.com/redis/go-redis/v9"
	"github.com/satori/go.uuid"
	redisStorage "kantoku/impl/l0/cell/redis"
	"kantoku/l0/cell"
	"testing"
)

func TestRedisStorage(t *testing.T) {
	client := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})
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

	storage := redisStorage.NewStorage(client)

	t.Run("create", func(t *testing.T) {
		id, err := storage.Create(context.Background(), []byte("test data"))
		if err != nil {
			t.Fatalf("create failed: %v", err)
		}
		if _, err := uuid.FromString(id); err != nil {
			t.Fatalf("invalid uuid: %v", err)
		}
	})

	t.Run("get", func(t *testing.T) {
		id, _ := storage.Create(context.Background(), []byte("test data"))
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

	t.Run("set", func(t *testing.T) {
		id, _ := storage.Create(context.Background(), []byte("test data"))
		c := &cell.Cell{
			ID:   id,
			Data: []byte("new data"),
		}
		if err := storage.Set(context.Background(), c); err != nil {
			t.Fatalf("set failed: %v", err)
		}
		c, err := storage.Get(context.Background(), id)
		if err != nil {
			t.Fatalf("get failed: %v", err)
		}
		if string(c.Data) != "new data" {
			t.Fatalf("incorrect data, expected %q but got %q", "new data", string(c.Data))
		}
	})

	t.Run("delete", func(t *testing.T) {
		id, _ := storage.Create(context.Background(), []byte("test data"))
		if err := storage.Delete(context.Background(), id); err != nil {
			t.Fatalf("delete failed: %v", err)
		}
		_, err := storage.Get(context.Background(), id)
		if err == nil {
			t.Fatalf("c should be deleted")
		}
	})
}

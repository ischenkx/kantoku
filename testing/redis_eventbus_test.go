package testing

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/redis/go-redis/v9"
	redisEB "hayaku/impl/l0/event/redis"
	"hayaku/l0/event"
	"io"
	"strconv"
	"testing"
	"time"
)

type TestCodec struct{}

func (t TestCodec) Encode(e event.Event) ([]byte, error) {
	data, err := json.Marshal(e)
	if err != nil {
		return nil, err
	}
	return data, nil
}

func (t TestCodec) Decode(reader io.Reader) (event.Event, error) {
	var e event.Event
	err := json.NewDecoder(reader).Decode(&e)
	if err != nil {
		return event.Event{}, err
	}
	return e, nil
}

func TestRedisEventBus(t *testing.T) {
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

	bus := redisEB.NewBus(client, TestCodec{})

	t.Run("publish", func(t *testing.T) {
		err := bus.Publish(context.Background(), event.Event{
			Name:  "test_publish_1",
			Topic: "test_topic",
			Data:  []byte("test data"),
		})

		if err != nil {
			t.Fatalf("publishing event failed: %v", err)
		}
	})

	t.Run("listen", func(t *testing.T) {
		ctx := context.Background()
		defer ctx.Done()

		_, err := bus.Listen(ctx, "test_topic")
		if err != nil {
			t.Fatalf("listening failed: %v", err)
		}
	})

	t.Run("listen and publish", func(t *testing.T) {
		ctx := context.Background()
		defer ctx.Done()

		channel, err := bus.Listen(ctx, "test_topic")
		if err != nil {
			t.Fatalf("listening failed: %v", err)
		}
		eSent := event.Event{
			Name:  "test_publish_2",
			Topic: "test_topic",
			Data:  []byte("test data"),
		}

		go func() {
			time.Sleep(1 * time.Second)
			err2 := bus.Publish(ctx, eSent)
			if err2 != nil {
				t.Errorf("publishing event failed: %v", err2)
				return
			}
		}()

		select {
		case eReceived := <-channel:
			if eSent.Name != eReceived.Name {
				t.Fatalf("received event with different name: sent='%v' received='%v'", eSent.Name, eReceived.Name)
			}
			if eSent.Topic != eReceived.Topic {
				t.Fatalf("received event with different topic: sent='%v' received='%v'", eSent.Topic, eReceived.Topic)
			}
			if string(eSent.Data) != string(eReceived.Data) {
				t.Fatalf("received event with different data: sent='%v' received='%v'",
					string(eSent.Data),
					string(eReceived.Data))
			}
		case <-time.After(3 * time.Second):
			t.Fatalf("did not receive event")
		}
	})

	t.Run("many events", func(t *testing.T) {
		ctx := context.Background()
		defer ctx.Done()

		channel, err := bus.Listen(ctx, "test_topic")
		if err != nil {
			t.Fatalf("listening failed: %v", err)
		}

		const eventsCount = 1000

		go func() {
			time.Sleep(1 * time.Second)
			for i := 0; i < eventsCount; i++ {
				err2 := bus.Publish(ctx, event.Event{
					Name:  "test_publish_3",
					Topic: "test_topic",
					Data:  []byte(fmt.Sprint(i)),
				})
				if err2 != nil {
					t.Errorf("publishing event failed: %v", err)
					return
				}
			}
		}()
		received := map[int]struct{}{}

	loop:
		for {
			select {
			case eReceived := <-channel:
				if eReceived.Name != "test_publish_3" {
					t.Fatalf("received event with wrong name: '%v'", eReceived.Name)
				}
				if value, err2 := strconv.Atoi(string(eReceived.Data)); err2 != nil {
					t.Fatalf("event with data that cannot be converted to int: %v", string(eReceived.Data))
				} else {
					received[value] = struct{}{}
				}
				if len(received) == eventsCount {
					break loop
				}
			case <-time.After(3 * time.Second):
				break loop
			}
		}

		for i := 0; i < eventsCount; i++ {
			if _, ok := received[i]; !ok {
				t.Fatalf("lost event number %v", i)
			}
		}
	})

	t.Run("two topics", func(t *testing.T) {
		ctx := context.Background()
		defer ctx.Done()

		channel1, err := bus.Listen(ctx, "test_topic_1")
		if err != nil {
			t.Fatalf("listening failed: %v", err)
		}
		channel2, err := bus.Listen(ctx, "test_topic_2")
		if err != nil {
			t.Fatalf("listening failed: %v", err)
		}
		go func() {
			time.Sleep(1 * time.Second)
			err2 := bus.Publish(ctx, event.Event{
				Name:  "1",
				Topic: "test_topic_1",
				Data:  []byte{},
			})
			if err2 != nil {
				t.Errorf("publishing event failed: %v", err2)
				return
			}
			err2 = bus.Publish(ctx, event.Event{
				Name:  "2",
				Topic: "test_topic_2",
				Data:  []byte{},
			})
			if err2 != nil {
				t.Errorf("publishing event failed: %v", err2)
				return
			}
		}()

		receivedCount := 0
	loop:
		for {
			select {
			case e1 := <-channel1:
				if e1.Name != "1" {
					t.Fatalf("received event with wrong name: %v", e1.Name)
				} else {
					receivedCount++
					if receivedCount == 2 {
						break loop
					}
				}
			case e2 := <-channel2:
				if e2.Name != "2" {
					t.Fatalf("received event with wrong name: %v", e2.Name)
				} else {
					receivedCount++
					if receivedCount == 2 {
						break loop
					}
				}
			case <-time.After(3 * time.Second):
				t.Fatalf("did not receive event")
			}
		}
	})
}

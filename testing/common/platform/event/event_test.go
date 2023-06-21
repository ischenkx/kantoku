package event

import (
	"context"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	"github.com/samber/lo"
	"github.com/stretchr/testify/assert"
	"kantoku/impl/common/codec/jsoncodec"
	redivent "kantoku/impl/platform/event/redis"
	"kantoku/kernel/platform"
	"log"
	"math/rand"
	"testing"
	"time"
)

type Item struct {
	Data string
	Name string
}

func newRedisEvents(ctx context.Context) *redivent.Broker {
	client := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379", // Redis server address
		Password: "",               // Redis server password (leave empty if not set)
		DB:       0,                // Redis database index
	})

	if cmd := client.Ping(ctx); cmd.Err() != nil {
		panic("failed to ping the redis client: " + cmd.Err().Error())
	}

	return redivent.New(jsoncodec.New[platform.Event](), client)
}

func TestEvents(t *testing.T) {
	ctx := context.Background()
	implementations := map[string]platform.Broker{
		"redis": newRedisEvents(ctx),
	}

	for label, impl := range implementations {
		t.Run(label, func(t *testing.T) {
			listener := impl.Listen()
			defer listener.Close(ctx)

			channel, err := listener.Incoming(ctx)
			if err != nil {
				t.Fatal("failed to make a channel:", err)
			}

			totalTopics := 10 + rand.Intn(100)
			topics := lo.Times(totalTopics, func(index int) string {
				return lo.RandomString(100, lo.AlphanumericCharset)
			})

			for _, topic := range topics {
				if err := listener.Subscribe(ctx, topic); err != nil {
					t.Fatal("failed to subscribe to a topic:", err)
				}
			}

			totalEvents := 10 + rand.Intn(500)
			events := lo.Times(totalEvents, func(index int) platform.Event {
				return platform.Event{
					Name:  uuid.New().String(),
					Data:  []byte(uuid.New().String()),
					Topic: topics[rand.Intn(totalTopics)],
				}
			})

			for _, event := range events {
				if err := impl.Publish(ctx, event); err != nil {
					t.Fatal("failed to publish:", err)
				}
			}

			for range events {
				select {
				case incoming := <-channel:
					assert.Contains(t, events, incoming)
				case <-time.After(time.Second * 5):
					t.Fatal("failed to receive an incoming message within five seconds...")
				}
			}

			t.Log("testing partial unsubscribing")
			unsubscribedTopics := lo.Subset(topics, rand.Intn(totalTopics/2), uint(rand.Intn(totalTopics/2)))
			if err := listener.Unsubscribe(ctx, unsubscribedTopics...); err != nil {
				t.Fatal("failed to unsubscribe:", err)
			}

			subscribedTopics, _ := lo.Difference(topics, unsubscribedTopics)
			for _, topic := range unsubscribedTopics {
				err := impl.Publish(ctx, platform.Event{
					Name:  "UNSUBSCRIBED",
					Topic: topic,
					Data:  nil,
				})

				if err != nil {
					t.Fatal("failed to published to an unsubscribed topic:", err)
				}
			}

			select {
			case event := <-channel:
				if event.Name == "UNSUBSCRIBED" {
					t.Fatal("received event from an unsubscribed topic:")
				} else {
					t.Fatal("received an unexpected event:", event.Name)
				}
			case <-time.After(time.Second * 5):
			}

			t.Log("testing partial subscribing")
			for _, topic := range subscribedTopics {
				err := impl.Publish(ctx, platform.Event{
					Name:  "SUBSCRIBED",
					Topic: topic,
					Data:  nil,
				})

				if err != nil {
					t.Fatal("failed to published to an unsubscribed topic:", err)
				}
			}

			for i := 0; i < len(subscribedTopics); i++ {
				select {
				case event := <-channel:
					if event.Name != "SUBSCRIBED" {
						t.Fatal("received an unexpected event:", event.Name)
					}
				case <-time.After(time.Second * 5):
					t.Fatal("expected a \"SUBSCRIBED\" event but received nothing")
				}
			}

			t.Log("testing unsubscribe all")
			if err := listener.UnsubscribeAll(ctx); err != nil {
				log.Fatal("failed to unsubscribe all:", err)
			}

			testMessages := rand.Intn(100) + 5
			for i := 0; i < testMessages; i++ {
				err := impl.Publish(ctx, platform.Event{
					Name:  "UNSUBSCRIBED_FROM_ALL",
					Topic: topics[rand.Intn(totalTopics)],
				})
				if err != nil {
					t.Fatal("failed to publish a test message:", err)
				}
			}

			select {
			case event := <-channel:
				if event.Name == "UNSUBSCRIBED_FROM_ALL" {
					t.Fatal("received event from an unsubscribed topic")
				} else {
					t.Fatal("received an unexpected event:", event.Name)
				}
			case <-time.After(time.Second * 5):
			}

			t.Log("Done!")
		})
	}
}

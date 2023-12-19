package event

import (
	"context"
	redisEvents "github.com/ischenkx/kantoku/pkg/impl/kernel/event/redis"
	event2 "github.com/ischenkx/kantoku/pkg/system/kernel/event"
	"github.com/redis/go-redis/v9"
	"github.com/samber/lo"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

var TimeoutPerCollectedItem = time.Second * 2

func newRedisEvents(ctx context.Context) *redisEvents.Bus {
	client := redis.NewUniversalClient(&redis.UniversalOptions{
		Addrs: []string{"172.23.146.206:6379"},
		//Addrs: []string{":6379"},
		DB: 1,
	})

	if err := client.Ping(ctx).Err(); err != nil {
		panic(err)
	}

	if _, err := client.FlushDB(ctx).Result(); err != nil {
		panic(err)
	}

	return redisEvents.New(client, redisEvents.StreamSettings{
		BatchSize:         64,
		ChannelBufferSize: 64,
		Consumer:          "tester-0",
	})
}

func TestBus(t *testing.T) {
	ctx := context.Background()

	implementations := map[string]event2.Bus{
		"redis": newRedisEvents(ctx),
	}

	for name, impl := range implementations {
		t.Run(name, func(t *testing.T) {
			ImplTest(ctx, t, impl)
		})
	}
}

func ImplTest(ctx context.Context, t *testing.T, bus event2.Bus) {
	events := []event2.Event{
		event2.New("a", []byte("1")),
		event2.New("a", []byte("2")),
		event2.New("b", []byte("3")),
		event2.New("b", []byte("4")),
		event2.New("b", []byte("5")),
		event2.New("b", []byte("6")),
		event2.New("b", []byte("7")),
		event2.New("c", []byte("8")),
		event2.New("c", []byte("9")),
		event2.New("c", []byte("10")),
		event2.New("c", []byte("11")),
		event2.New("d", []byte("12")),
		event2.New("d", []byte("13")),
		event2.New("d", []byte("13")),
		event2.New("d", []byte("13")),
	}
	queues := []event2.Queue{
		{
			Name: "q1",
			Topics: []string{
				"a", "b",
			},
		},
		{
			Name: "q2",
			Topics: []string{
				"a",
			},
		},
		{
			Name: "q3",
			Topics: []string{
				"b", "c",
			},
		},
		{
			Name: "q4",
			Topics: []string{
				"b", "d",
			},
		},
	}

	for _, ev := range events {
		if err := bus.Publish(ctx, ev); err != nil {
			t.Fatalf("failed to publish an event: %s", err)
		}
	}

	for _, queue := range queues {
		t.Logf("testing queue (name = '%s', topics = '%v')\n",
			queue.Name,
			queue.Topics)
		collectedEvents, err := collectFromQueue(ctx, queue, bus)
		if err != nil {
			t.Fatal("failed to collect:", err)
		}

		expectedEvents := filterEventsByTopics(events, queue.Topics)

		assert.ElementsMatch(t, expectedEvents, collectedEvents)
	}
}

func filterEventsByTopics(events []event2.Event, topics []string) []event2.Event {
	return lo.Filter(events, func(ev event2.Event, _ int) bool {
		return lo.Contains(topics, ev.Topic)
	})
}

func collectFromQueue(ctx context.Context, queue event2.Queue, bus event2.Bus) ([]event2.Event, error) {
	var collectedEvents []event2.Event

	channel, err := bus.Consume(ctx, queue)
	if err != nil {
		return nil, err
	}

collector:
	for {
		select {
		case <-time.After(TimeoutPerCollectedItem):
			break collector
		case ev := <-channel:
			collectedEvents = append(collectedEvents, ev)
		}
	}

	return collectedEvents, nil
}

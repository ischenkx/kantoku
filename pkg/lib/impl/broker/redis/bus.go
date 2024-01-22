package redis

import (
	"context"
	"errors"
	"fmt"
	event2 "github.com/ischenkx/kantoku/pkg/system/kernel/event"
	"github.com/redis/go-redis/v9"
	"github.com/samber/lo"
	"log"
	"strconv"
	"strings"
)

type StreamSettings struct {
	BatchSize         int
	ChannelBufferSize int
	Consumer          string
}

type Bus struct {
	settings StreamSettings
	client   redis.UniversalClient
}

func New(client redis.UniversalClient, settings StreamSettings) *Bus {
	return &Bus{
		settings: settings,
		client:   client,
	}
}

func (bus *Bus) Consume(ctx context.Context, queue event2.Queue) (<-chan event2.Event, error) {
	if queue.Name == "" {
		return nil, errors.New("empty queue name")
	}

	channel := make(chan event2.Event, bus.settings.ChannelBufferSize)

	for _, topic := range queue.Topics {
		if err := bus.initializeGroup(ctx, topic, queue.Name); err != nil {
			return nil, fmt.Errorf("failed to initialize a group (topic: %s, group: %s): %w", topic, queue.Name, err)
		}
	}

	go bus.groupReadToChannel(ctx, queue.Name, queue.Topics, channel)

	return channel, nil
}

func (bus *Bus) initializeGroup(ctx context.Context, stream, group string) error {
	err := bus.client.XGroupCreateMkStream(ctx, stream, group, "0").Err()
	if err != nil {
		if strings.Contains(err.Error(), "BUSYGROUP") {
			err = nil
		}
	}

	return err
}

func (bus *Bus) groupReadToChannel(ctx context.Context, group string, streams []string, destination chan<- event2.Event) {
	args := redis.XReadGroupArgs{
		Group:    group,
		Consumer: bus.settings.Consumer,
		Streams: append(streams, lo.RepeatBy(len(streams), func(_ int) string {
			return ">"
		})...),
		Count: int64(bus.settings.BatchSize),
		Block: 0,
		NoAck: true,
	}

poller:
	for {
		select {
		case <-ctx.Done():
			break poller
		default:
		}

		streams, err := bus.client.XReadGroup(ctx, &args).Result()
		if err != nil {
			log.Println("failed to read the group:", err)
			continue poller
		}

		for _, stream := range streams {
			for _, message := range stream.Messages {
				ev, err := bus.message2event(message.Values)
				if err != nil {
					log.Println("failed to parse the event:", err)
					continue
				}

				select {
				case <-ctx.Done():
					log.Println("some messages might be dropped due to the context cancellation")
					break poller
				case destination <- ev:
				}
			}
		}
	}
}

//func (bus *Bus) readToChannel(ctx context.Context, streams []string, destination chan<- event2.Item) {
//	args := redis.XReadArgs{
//		Streams: append(streams, lo.RepeatBy(len(streams), func(_ int) string {
//			return ">"
//		})...),
//		Count: int64(bus.settings.BatchSize),
//		Block: 0,
//	}
//
//poller:
//	for {
//		select {
//		case <-ctx.Done():
//			break poller
//		default:
//		}
//
//		streams, err := bus.client.XReadGroup(ctx, &args).Result()
//		if err != nil {
//			log.Println("failed to read the group:", err)
//			continue poller
//		}
//
//		for _, stream := range streams {
//			for _, message := range stream.Messages {
//				ev, err := bus.message2event(message.Values)
//				if err != nil {
//					log.Println("failed to parse the event:", err)
//					continue
//				}
//
//				select {
//				case <-ctx.Done():
//					log.Println("some messages might be dropped due to the context cancellation")
//					break poller
//				case destination <- ev:
//				}
//			}
//		}
//	}
//}

func (bus *Bus) event2message(ev event2.Event) map[string]any {
	return map[string]any{
		"data":      ev.Data,
		"topic":     ev.Topic,
		"timestamp": ev.Timestamp,
	}
}

func (bus *Bus) message2event(message map[string]any) (event2.Event, error) {
	var ev event2.Event

	if data, ok := message["data"]; ok {
		ev.Data = []byte(data.(string))
	} else {
		return ev, errors.New("no data present in the message")
	}

	if data, ok := message["topic"]; ok {
		ev.Topic = data.(string)
	} else {
		return ev, errors.New("no topic present in the message")
	}

	if rawTimestamp, ok := message["timestamp"]; ok {
		timestamp, err := strconv.ParseInt(rawTimestamp.(string), 0, 64)
		if err != nil {
			return ev, err
		}

		ev.Timestamp = timestamp
	} else {
		return ev, errors.New("no data present in message")
	}

	return ev, nil
}

func (bus *Bus) Publish(ctx context.Context, event event2.Event) error {
	_, err := bus.client.
		XAdd(ctx, &redis.XAddArgs{
			Stream: event.Topic,
			ID:     "*",
			Values: bus.event2message(event),
		}).
		Result()
	if err != nil {
		return fmt.Errorf("failed to xadd: %w", err)
	}

	return nil
}

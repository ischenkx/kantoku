package broker

import (
	"context"
)

type ConsumerInitializationPolicy string

var (
	OldestOffset ConsumerInitializationPolicy = "oldest"
	NewestOffset ConsumerInitializationPolicy = "newest"
)

type ConsumerSettings struct {
	Group                string
	InitializationPolicy ConsumerInitializationPolicy
}

type Consumer[Item any] interface {
	Consume(ctx context.Context, topics []string, settings ConsumerSettings) (<-chan Message[Item], error)
}

type Publisher[Item any] interface {
	Publish(ctx context.Context, topic string, item Item) error
}

type Broker[Item any] interface {
	Consumer[Item]
	Publisher[Item]
}

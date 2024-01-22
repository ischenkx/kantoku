package broker

import "context"

type TopicsInfo struct {
	Group  string
	Topics []string
}

type Consumer[Item any] interface {
	Consume(ctx context.Context, info TopicsInfo) (<-chan Message[Item], error)
}

type Publisher[Item any] interface {
	Publish(ctx context.Context, topic string, item Item) error
}

type Broker[Item any] interface {
	Consumer[Item]
	Publisher[Item]
}

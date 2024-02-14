package broker

import (
	"context"
	"github.com/ischenkx/kantoku/pkg/common/transport/queue"
)

type TopicsInfo struct {
	Group  string
	Topics []string
}

type Consumer[Item any] interface {
	Consume(ctx context.Context, info TopicsInfo) (<-chan queue.Message[Item], error)
}

type Publisher[Item any] interface {
	Publish(ctx context.Context, topic string, item Item) error
}

type Broker[Item any] interface {
	Consumer[Item]
	Publisher[Item]
}

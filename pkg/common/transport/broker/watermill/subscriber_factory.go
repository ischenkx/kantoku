package watermill

import (
	"context"
	"github.com/ThreeDotsLabs/watermill/message"
)

type SubscriberFactory interface {
	New(ctx context.Context, consumerGroup string) (message.Subscriber, error)
}

type FunctionalSubscriberFactory func(ctx context.Context, consumerGroup string) (message.Subscriber, error)

func (f FunctionalSubscriberFactory) New(ctx context.Context, consumerGroup string) (message.Subscriber, error) {
	return f(ctx, consumerGroup)
}

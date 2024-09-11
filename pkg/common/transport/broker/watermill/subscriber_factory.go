package watermill

import (
	"context"
	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/ischenkx/kantoku/pkg/common/transport/broker"
)

type SubscriberFactory interface {
	New(ctx context.Context, settings broker.ConsumerSettings) (message.Subscriber, error)
}

type FunctionalSubscriberFactory func(
	ctx context.Context,
	settings broker.ConsumerSettings,
) (message.Subscriber, error)

func (f FunctionalSubscriberFactory) New(ctx context.Context, settings broker.ConsumerSettings) (message.Subscriber, error) {
	return f(ctx, settings)
}

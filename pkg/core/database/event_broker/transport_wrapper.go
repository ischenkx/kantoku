package eventbroker

import (
	"context"
	"github.com/ischenkx/kantoku/pkg/common/transport/broker"
	"github.com/ischenkx/kantoku/pkg/core"
)

var _ core.Broker = (*CommonBrokerWrapper)(nil)

type CommonBrokerWrapper struct {
	broker broker.Broker[core.Event]
}

func WrapCommonBroker(b broker.Broker[core.Event]) *CommonBrokerWrapper {
	return &CommonBrokerWrapper{broker: b}
}

func (b *CommonBrokerWrapper) Send(ctx context.Context, event core.Event) error {
	return b.broker.Publish(ctx, event.Topic, event)
}

func (b *CommonBrokerWrapper) Consume(ctx context.Context, events []string, settings broker.ConsumerSettings) (<-chan core.BrokerEvent, error) {
	return b.broker.Consume(ctx, events, settings)
}

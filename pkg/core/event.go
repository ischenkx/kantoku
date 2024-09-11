package core

import (
	"context"
	"github.com/google/uuid"
	"github.com/ischenkx/kantoku/pkg/common/transport/broker"
	"time"
)

type Event struct {
	ID        string
	Data      []byte
	Topic     string
	Timestamp int64
}

func NewEvent(topic string, data []byte) Event {
	return Event{
		ID:        uuid.New().String(),
		Data:      data,
		Topic:     topic,
		Timestamp: time.Now().UnixNano(),
	}
}

type BrokerEvent = broker.Message[Event]

type Broker interface {
	Send(ctx context.Context, event Event) error
	Consume(ctx context.Context, events []string, consumerSettings broker.ConsumerSettings) (<-chan BrokerEvent, error)
}

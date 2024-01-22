package event

import (
	"github.com/google/uuid"
	"time"
)

type Event struct {
	ID        string
	Data      []byte
	Topic     string
	Timestamp int64
}

func New(topic string, data []byte) Event {
	return Event{
		ID:        uuid.New().String(),
		Data:      data,
		Topic:     topic,
		Timestamp: time.Now().UnixNano(),
	}
}

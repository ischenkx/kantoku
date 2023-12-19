package event

import "time"

type Event struct {
	Data      []byte
	Topic     string
	Timestamp int64
}

func New(topic string, data []byte) Event {
	return Event{
		Data:      data,
		Topic:     topic,
		Timestamp: time.Now().Unix(),
	}
}

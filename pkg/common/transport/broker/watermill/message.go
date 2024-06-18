package watermill

import "github.com/ThreeDotsLabs/watermill/message"

type Message[Item any] struct {
	item  Item
	topic string
	raw   *message.Message
}

func (mes Message[Item]) Topic() string {
	return mes.topic
}

func (mes Message[Item]) Item() Item {
	return mes.item
}

func (mes Message[Item]) Ack() {
	mes.raw.Ack()
}

func (mes Message[Item]) Nack() {
	mes.raw.Nack()
}

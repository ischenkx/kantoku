package broker

type Message[Item any] interface {
	Item() Item
	Topic() string
	Ack()
	Nack()
}

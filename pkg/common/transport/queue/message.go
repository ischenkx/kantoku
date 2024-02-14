package queue

type Message[Item any] interface {
	Item() Item
	Ack()
	Nack()
}

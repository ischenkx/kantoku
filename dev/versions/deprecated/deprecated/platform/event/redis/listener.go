package redivent

import (
	"context"
	"github.com/redis/go-redis/v9"
	"kantoku/common/codec"
	"kantoku/job/platform"
	"log"
	"sync"
	"time"
)

type Listener struct {
	codec     codec.Codec[platform.Event, []byte]
	receivers map[chan<- platform.Event][]platform.Event
	pubsub    *redis.PubSub
	mu        sync.Mutex
}

func NewListener(codec codec.Codec[platform.Event, []byte], pubsub *redis.PubSub) *Listener {
	listener := &Listener{
		codec:     codec,
		receivers: map[chan<- platform.Event][]platform.Event{},
		pubsub:    pubsub,
	}
	go listener.runFlusher()
	go listener.runBroadcaster()
	return listener
}

func (listener *Listener) Subscribe(ctx context.Context, topics ...string) error {
	return listener.pubsub.Subscribe(ctx, topics...)
}

func (listener *Listener) Unsubscribe(ctx context.Context, topics ...string) error {
	if len(topics) == 0 {
		return nil
	}
	return listener.pubsub.Unsubscribe(ctx, topics...)
}

func (listener *Listener) UnsubscribeAll(ctx context.Context) error {
	return listener.pubsub.Unsubscribe(ctx)
}

func (listener *Listener) Incoming(ctx context.Context) (<-chan platform.Event, error) {
	return listener.initReceiver(ctx), nil
}

func (listener *Listener) Close(_ context.Context) error {
	listener.mu.Lock()
	defer listener.mu.Unlock()
	for receiver := range listener.receivers {
		close(receiver)
	}
	listener.receivers = nil
	return listener.pubsub.Close()
}

func (listener *Listener) initReceiver(ctx context.Context) <-chan platform.Event {
	listener.mu.Lock()
	defer listener.mu.Unlock()
	channel := make(chan platform.Event, 512)
	listener.receivers[channel] = []platform.Event{}

	go func(ctx context.Context, receiver chan<- platform.Event) {
		<-ctx.Done()
		listener.removeReceiver(receiver)
	}(ctx, channel)

	return channel
}

func (listener *Listener) removeReceiver(receiver chan<- platform.Event) {
	listener.mu.Lock()
	defer listener.mu.Unlock()
	close(receiver)
	delete(listener.receivers, receiver)
}

func (listener *Listener) runBroadcaster() {
	for message := range listener.pubsub.Channel() {
		event, err := listener.codec.Decode([]byte(message.Payload))
		if err != nil {
			log.Println("failed to decode an incoming message:", err)
		}

		listener.broadcast(event)
	}
}

func (listener *Listener) runFlusher() {
	ticker := time.NewTicker(time.Second * 3)
	defer ticker.Stop()

	for range ticker.C {
		listener.mu.Lock()
		if listener.receivers == nil {
			listener.mu.Unlock()
			break
		}
		for receiver := range listener.receivers {
			listener.unsafeFlush(receiver)
		}
		listener.mu.Unlock()
	}
}

func (listener *Listener) broadcast(event platform.Event) {
	listener.mu.Lock()
	defer listener.mu.Unlock()

	for receiver := range listener.receivers {
		listener.receivers[receiver] = append(listener.receivers[receiver], event)
		listener.unsafeFlush(receiver)
	}
}

func (listener *Listener) unsafeFlush(receiver chan<- platform.Event) {
	buffer := listener.receivers[receiver]
flusher:
	for i := 0; i < len(buffer); i++ {
		select {
		case receiver <- buffer[0]:
		default:
			break flusher
		}
		buffer = buffer[1:]
	}
	listener.receivers[receiver] = buffer
}

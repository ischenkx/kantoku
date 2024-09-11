package watermill

import (
	"context"
	"fmt"
	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/google/uuid"
	"github.com/ischenkx/kantoku/pkg/common/data/codec"
	"github.com/ischenkx/kantoku/pkg/common/transport/broker"
	"github.com/samber/lo"
	"log/slog"
)

type Broker[Item any] struct {
	Agent                     Agent
	ItemCodec                 codec.Codec[Item, []byte]
	Logger                    *slog.Logger
	ConsumerChannelBufferSize int
}

func (b Broker[Item]) Consume(ctx context.Context, topics []string, settings broker.ConsumerSettings) (<-chan broker.Message[Item], error) {
	subscriber, err := b.Agent.SubscriberFactory.New(ctx, settings)
	if err != nil {
		return nil, fmt.Errorf("failed to create a subscriber: %w", err)
	}

	channels := make([]<-chan *message.Message, 0, len(topics))
	for _, topic := range topics {
		ctx := context.WithValue(ctx, "topic", topic)
		channel, err := subscriber.Subscribe(ctx, topic)
		if err != nil {
			if closeErr := subscriber.Close(); closeErr != nil {
				b.Logger.Error("failed to close a subscriber",
					slog.String("error", closeErr.Error()))
			}

			b.Logger.Error("failed to subscribe",
				slog.String("error", err.Error()))

			return nil, fmt.Errorf("failed to subscribe: %w", err)
		}

		channels = append(channels, channel)
	}

	resultChannel := make(chan broker.Message[Item], b.ConsumerChannelBufferSize)

	mergedChannel := lo.FanIn[*message.Message](b.ConsumerChannelBufferSize, channels...)
	go func(ctx context.Context, from <-chan *message.Message, to chan<- broker.Message[Item]) {
		defer subscriber.Close()

		for {
			select {
			case <-ctx.Done():
				return
			case mes := <-from:
				item, err := b.ItemCodec.Decode(mes.Payload)
				if err != nil {
					b.Logger.Error("failed to decode item",
						slog.String("error", err.Error()))
					continue
				}

				rawTopic := mes.Context().Value("topic")
				topic, ok := rawTopic.(string)
				if !ok {
					b.Logger.Error("failed to extract topic from message",
						slog.Any("raw_topic", rawTopic))
					continue
				}

				brokerMessage := Message[Item]{
					item:  item,
					topic: topic,
					raw:   mes,
				}

				select {
				case <-ctx.Done():
					return
				case to <- brokerMessage:
				}
			}
		}
	}(ctx, mergedChannel, resultChannel)

	return resultChannel, nil
}

func (b Broker[Item]) Publish(_ context.Context, topic string, item Item) error {
	payload, err := b.ItemCodec.Encode(item)
	if err != nil {
		return fmt.Errorf("failed to encode item: %w", err)
	}

	if err := b.Agent.Publisher.Publish(topic, message.NewMessage(uuid.New().String(), payload)); err != nil {
		return err
	}

	return nil
}

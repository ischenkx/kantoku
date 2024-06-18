package watermill

import (
	"context"
	"fmt"
	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill-nats/v2/pkg/nats"
	"github.com/ThreeDotsLabs/watermill-redisstream/pkg/redisstream"
	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/redis/go-redis/v9"
	"log/slog"
)

// Redis

func Redis(
	client redis.UniversalClient,
	subscriberConfigTemplate redisstream.SubscriberConfig,
	publisherConfigTemplate redisstream.PublisherConfig,
) (Agent, error) {
	subscriberFactory := FunctionalSubscriberFactory(
		func(ctx context.Context, consumerGroup string) (message.Subscriber, error) {
			configTemplate := subscriberConfigTemplate
			configTemplate.ConsumerGroup = consumerGroup
			configTemplate.Client = client
			return redisstream.NewSubscriber(
				configTemplate,
				watermill.NewSlogLogger(slog.Default()),
			)
		},
	)

	publisherConfigTemplate.Client = client
	publisher, err := redisstream.NewPublisher(
		publisherConfigTemplate,
		watermill.NewSlogLogger(slog.Default()),
	)
	if err != nil {
		return Agent{}, fmt.Errorf("failed to create a publisher: %w", err)
	}

	return Agent{SubscriberFactory: subscriberFactory, Publisher: publisher}, nil
}

// NATS

func Nats(
	url string,
	subscriberConfigTemplate nats.SubscriberConfig,
	publisherConfigTemplate nats.PublisherConfig,
	logger *slog.Logger,
) (Agent, error) {
	subscriberFactory := FunctionalSubscriberFactory(
		func(ctx context.Context, consumerGroup string) (message.Subscriber, error) {
			configTemplate := subscriberConfigTemplate

			configTemplate.URL = url
			configTemplate.QueueGroupPrefix = consumerGroup
			configTemplate.JetStream.DurablePrefix = consumerGroup
			return nats.NewSubscriber(
				configTemplate,
				watermill.NewSlogLogger(logger),
			)
		},
	)

	publisherConfigTemplate.URL = url
	publisher, err := nats.NewPublisher(
		publisherConfigTemplate,
		watermill.NewSlogLogger(logger),
	)
	if err != nil {
		return Agent{}, fmt.Errorf("failed to create a publisher: %w", err)
	}

	return Agent{SubscriberFactory: subscriberFactory, Publisher: publisher}, nil
}

package watermill

import (
	"context"
	"fmt"
	"github.com/Shopify/sarama"
	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill-kafka/v2/pkg/kafka"
	"github.com/ThreeDotsLabs/watermill-nats/v2/pkg/nats"
	"github.com/ThreeDotsLabs/watermill-redisstream/pkg/redisstream"
	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/ischenkx/kantoku/pkg/common/transport/broker"
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
		func(ctx context.Context, settings broker.ConsumerSettings) (message.Subscriber, error) {
			configTemplate := subscriberConfigTemplate
			configTemplate.ConsumerGroup = settings.Group
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
		func(ctx context.Context, settings broker.ConsumerSettings) (message.Subscriber, error) {
			configTemplate := subscriberConfigTemplate

			configTemplate.URL = url
			configTemplate.QueueGroupPrefix = settings.Group
			configTemplate.JetStream.DurablePrefix = settings.Group

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

// Kafka

func Kafka(
	brokers []string,
	subscriberConfigTemplate kafka.SubscriberConfig,
	publisherConfigTemplate kafka.PublisherConfig,
	logger *slog.Logger,
) (Agent, error) {
	subscriberFactory := FunctionalSubscriberFactory(
		func(ctx context.Context, settings broker.ConsumerSettings) (message.Subscriber, error) {
			configTemplate := subscriberConfigTemplate

			saramaConfig := kafka.DefaultSaramaSubscriberConfig()
			saramaConfig.Consumer.Offsets.Initial = sarama.OffsetOldest
			if settings.InitializationPolicy == broker.NewestOffset {
				saramaConfig.Consumer.Offsets.Initial = sarama.OffsetNewest
			}

			configTemplate.OverwriteSaramaConfig = saramaConfig
			configTemplate.ConsumerGroup = settings.Group
			configTemplate.Brokers = brokers

			return kafka.NewSubscriber(configTemplate, watermill.NewSlogLogger(logger))
		},
	)

	publisherConfigTemplate.Brokers = brokers

	publisher, err := kafka.NewPublisher(
		publisherConfigTemplate,
		watermill.NewSlogLogger(logger),
	)
	if err != nil {
		return Agent{}, fmt.Errorf("failed to create a publisher: %w", err)
	}

	return Agent{SubscriberFactory: subscriberFactory, Publisher: publisher}, nil
}

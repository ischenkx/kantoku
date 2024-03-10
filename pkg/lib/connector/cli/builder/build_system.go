package builder

import (
	"context"
	"fmt"
	"github.com/ThreeDotsLabs/watermill-nats/v2/pkg/nats"
	"github.com/ischenkx/kantoku/pkg/common/data/codec"
	"github.com/ischenkx/kantoku/pkg/common/data/record"
	"github.com/ischenkx/kantoku/pkg/common/transport/broker"
	"github.com/ischenkx/kantoku/pkg/common/transport/queue"
	"github.com/ischenkx/kantoku/pkg/core/event"
	"github.com/ischenkx/kantoku/pkg/core/resource"
	"github.com/ischenkx/kantoku/pkg/core/system"
	"github.com/ischenkx/kantoku/pkg/core/task"
	"github.com/ischenkx/kantoku/pkg/lib/connector/cli/config"
	"github.com/ischenkx/kantoku/pkg/lib/connector/cli/errx"
	"github.com/ischenkx/kantoku/pkg/lib/impl/broker/watermill"
	redisResources "github.com/ischenkx/kantoku/pkg/lib/impl/core/resource/redis"
	mongorec "github.com/ischenkx/kantoku/pkg/lib/impl/data/record/mongo"
	"github.com/ischenkx/kantoku/pkg/lib/resources"
	nc "github.com/nats-io/nats.go"
	"log/slog"
	"os"
	"time"
)

func (builder *Builder) BuildSystem(ctx context.Context, cfg config.SystemConfig) (system.AbstractSystem, error) {
	sys := &system.System{}

	logger := newLogger(os.Stdout).With("service", "system")
	ctx = withLogger(ctx, logger)

	ctx = withSystem(ctx, sys)

	tasks, err := builder.BuildTasks(ctx, cfg.Tasks)
	if err != nil {
		return system.System{}, errx.FailedToBuild("tasks", err)
	}

	resourceStorage, err := builder.BuildResources(ctx, cfg.Resources)
	if err != nil {
		return system.System{}, errx.FailedToBuild("resources", err)
	}

	events, err := builder.BuildEvents(ctx, cfg.Events)
	if err != nil {
		return system.System{}, errx.FailedToBuild("events", err)
	}

	*sys = system.System{
		Events_:    events,
		Resources_: resourceStorage,
		Tasks_:     tasks,
		Logger:     logger,
	}

	return sys, nil
}

func (builder *Builder) BuildEvents(ctx context.Context, cfg config.DynamicConfig) (*event.Broker, error) {
	var eventsConfig struct {
		Broker config.DynamicConfig
	}
	if err := cfg.Bind(&eventsConfig); err != nil {
		return nil, errx.FailedToBind(err)
	}

	b, err := builder.buildEventBroker(ctx, eventsConfig.Broker)
	if err != nil {
		return nil, errx.FailedToBuild("event broker", err)
	}

	return event.NewBroker(b), nil
}

func (builder *Builder) buildEventBroker(ctx context.Context, cfg config.DynamicConfig) (broker.Broker[event.Event], error) {
	switch cfg.Kind() {
	case "nats":
		var natsConfig struct {
			URI string
		}
		if err := cfg.Bind(&natsConfig); err != nil {
			return nil, errx.FailedToBind(err)
		}

		natsOptions := []nc.Option{
			nc.RetryOnFailedConnect(true),
			nc.Timeout(30 * time.Second),
			nc.ReconnectWait(1 * time.Second),
		}
		subscribeOptions := []nc.SubOpt{
			nc.DeliverAll(),
			nc.AckExplicit(),
		}

		jsConfig := nats.JetStreamConfig{
			Disabled:         false,
			AutoProvision:    true,
			ConnectOptions:   nil,
			SubscribeOptions: subscribeOptions,
			PublishOptions:   nil,
			TrackMsgId:       false,
			AckAsync:         true,
			//DurablePrefix:    "kantoku",
			//DurableCalculator: func(prefix string, topic string) string {
			//	return fmt.Sprintf("%s:%s", prefix, topic)
			//},
		}

		subscriberConfig := nats.SubscriberConfig{
			//SubscribersCount: 1,
			//CloseTimeout:      time.Second * 30,
			//AckWaitTimeout:    time.Second * 30,
			//SubscribeTimeout:  time.Second * 30,
			NatsOptions:       natsOptions,
			Unmarshaler:       nil,
			SubjectCalculator: nil,
			NakDelay:          nil,
			JetStream:         jsConfig,
		}

		publishedConfig := nats.PublisherConfig{
			NatsOptions: natsOptions,
			JetStream:   jsConfig,
		}

		agent, err := watermill.Nats(
			natsConfig.URI,
			subscriberConfig,
			publishedConfig,
			extractLogger(ctx, slog.Default()),
		)
		if err != nil {
			return nil, fmt.Errorf("failed to connect to nats: %w", err)
		}

		b := watermill.Broker[event.Event]{
			Agent:                     agent,
			ItemCodec:                 codec.JSON[event.Event](),
			Logger:                    extractLogger(ctx, slog.Default()),
			ConsumerChannelBufferSize: 1024,
		}

		return b, nil
	default:
		return nil, errx.UnsupportedKind(cfg.Kind())
	}
}

func (builder *Builder) BuildResources(ctx context.Context, cfg config.DynamicConfig) (resource.Storage, error) {
	var resourcesConfig struct {
		Storage   config.DynamicConfig
		Observers []config.DynamicConfig
	}
	if err := cfg.Bind(&resourcesConfig); err != nil {
		return nil, errx.FailedToBind(err)
	}

	storage, err := builder.buildResourcesStorage(ctx, resourcesConfig.Storage)
	if err != nil {
		return nil, errx.FailedToBuild("resource_storage", err)
	}

	var observers []resources.Observer
	for _, observerConfig := range resourcesConfig.Observers {
		observer, err := builder.buildResourceObserver(ctx, observerConfig)
		if err != nil {
			return nil, errx.FailedToBuild(
				fmt.Sprintf("observer (%s)", observerConfig.Kind()),
				err)
		}

		observers = append(observers, observer)
	}

	return resources.Observe(storage, resources.MultiObserver(observers)), nil
}

func (builder *Builder) buildResourcesStorage(ctx context.Context, cfg config.DynamicConfig) (resource.Storage, error) {
	switch cfg.Kind() {
	case "redis":
		redisClient, err := builder.BuildRedis(ctx, cfg)
		if err != nil {
			return nil, errx.FailedToBuild("redis", err)
		}
		return redisResources.New(redisClient, codec.JSON[resource.Resource](), "resource"), nil
	default:
		return nil, errx.UnsupportedKind(cfg.Kind())
	}
}

func (builder *Builder) buildResourceObserver(ctx context.Context, cfg config.DynamicConfig) (resources.Observer, error) {
	switch cfg.Kind() {
	case "notifier":
		var notifierConfig struct {
			Topic string
		}
		if err := cfg.Bind(&notifierConfig); err != nil {
			return nil, errx.FailedToBind(err)
		}

		sys, ok := extractSystem(ctx)
		if !ok {
			return nil, fmt.Errorf("failed to extract system")
		}

		notifier := resources.Notifier{
			Logger: extractLogger(ctx, slog.Default()),
			Dst: queue.FunctionalPublisher[string]{
				Func: func(ctx context.Context, item string) error {
					if sys == nil {
						return fmt.Errorf("system is nil")
					}

					err := sys.Events().Send(ctx, event.New(notifierConfig.Topic, []byte(item)))
					if err != nil {
						return fmt.Errorf("failed to send an event: %w", err)
					}

					return nil
				},
			},
			Topic: notifierConfig.Topic,
		}

		return notifier, nil
	default:
		return nil, errx.UnsupportedKind(cfg.Kind())
	}
}

func (builder *Builder) BuildTasks(ctx context.Context, cfg config.DynamicConfig) (record.Storage[task.Task], error) {
	var tasksConfig struct {
		Storage config.DynamicConfig
	}
	if err := cfg.Bind(&tasksConfig); err != nil {
		return nil, errx.FailedToBind(err)
	}

	return builder.buildTasksStorage(ctx, tasksConfig.Storage)
}

func (builder *Builder) buildTasksStorage(ctx context.Context, cfg config.DynamicConfig) (record.Storage[task.Task], error) {
	switch cfg.Kind() {
	case "mongo":
		mongoInfo, err := builder.BuildMongo(ctx, cfg)
		if err != nil {
			return nil, errx.FailedToBuild("mongo", err)
		}
		return mongorec.New[task.Task](mongoInfo.GetCollection(), task.Codec{}), nil
	default:
		return nil, errx.UnsupportedKind(cfg.Kind())
	}
}

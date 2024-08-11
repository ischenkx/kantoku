package platform

import (
	"context"
	"fmt"
	"github.com/ThreeDotsLabs/watermill-nats/v2/pkg/nats"
	capi "github.com/hashicorp/consul/api"
	"github.com/ischenkx/kantoku/pkg/common/data/codec"
	"github.com/ischenkx/kantoku/pkg/common/data/storage"
	"github.com/ischenkx/kantoku/pkg/common/data/uid"
	"github.com/ischenkx/kantoku/pkg/common/dependency"
	batched2 "github.com/ischenkx/kantoku/pkg/common/dependency/postgres/batched"
	"github.com/ischenkx/kantoku/pkg/common/logging/prefixed"
	"github.com/ischenkx/kantoku/pkg/common/service"
	"github.com/ischenkx/kantoku/pkg/common/transport/broker/watermill"
	"github.com/ischenkx/kantoku/pkg/core/event"
	"github.com/ischenkx/kantoku/pkg/core/resource"
	redisResources "github.com/ischenkx/kantoku/pkg/core/resource/redis"
	"github.com/ischenkx/kantoku/pkg/core/services/executor"
	"github.com/ischenkx/kantoku/pkg/core/services/scheduler/dependencies"
	manager2 "github.com/ischenkx/kantoku/pkg/core/services/scheduler/dependencies/manager"
	resourceResolver2 "github.com/ischenkx/kantoku/pkg/core/services/scheduler/dependencies/manager/resolvers/resource_resolver"
	"github.com/ischenkx/kantoku/pkg/core/services/scheduler/dependencies/manager/task2group"
	"github.com/ischenkx/kantoku/pkg/core/services/status"
	"github.com/ischenkx/kantoku/pkg/core/system"
	"github.com/ischenkx/kantoku/pkg/core/task"
	"github.com/ischenkx/kantoku/pkg/lib/discovery"
	"github.com/ischenkx/kantoku/pkg/lib/discovery/consul"
	"github.com/ischenkx/kantoku/pkg/lib/platform/errx"
	"github.com/ischenkx/kantoku/pkg/lib/resources"
	"github.com/ischenkx/kantoku/pkg/lib/tasks/specification"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/lmittmann/tint"
	nc "github.com/nats-io/nats.go"
	"github.com/redis/go-redis/v9"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"io"
	"log/slog"
	"time"
)

func BuildSystem(ctx context.Context, logger *slog.Logger, config SystemConfig) (*system.System, error) {
	tasks, err := BuildTasks(ctx, logger, config.Tasks)
	if err != nil {
		return nil, errx.FailedToBuild("tasks", err)
	}

	events, err := BuildEvents(ctx, logger, config.Events)
	if err != nil {
		return nil, errx.FailedToBuild("events", err)
	}

	resources, err := BuildResources(ctx, events, logger, config.Resources)
	if err != nil {
		return nil, errx.FailedToBuild("resources", err)
	}

	return &system.System{
		Events_:    events,
		Resources_: resources,
		Tasks_:     tasks,
		Logger:     logger.With(slog.String("component", "system")),
	}, nil
}

func BuildTasks(ctx context.Context, logger *slog.Logger, config TasksConfig) (task.Storage, error) {
	storage, err := BuildTasksStorage(ctx, logger, config.Storage)
	if err != nil {
		return nil, errx.FailedToBuild("tasks_storage", err)
	}

	return storage, nil
}

func BuildTasksStorage(ctx context.Context, logger *slog.Logger, config TasksStorageConfig) (task.Storage, error) {
	switch config.Kind {
	case "mongo":
		client, err := buildMongo(ctx, config.URI)
		if err != nil {
			return nil, errx.FailedToBuild("mongo", err)
		}

		db, err := getOption[string](config.Options, "db")
		if err != nil {
			return nil, errx.FailedToBuild("mongo", err)
		}

		collection, err := getOption[string](config.Options, "collection")
		if err != nil {
			return nil, errx.FailedToBuild("mongo", err)
		}

		st := &storage.MongoStorage{
			Collection: client.Database(db).Collection(collection),
			Logger:     logger,
		}

		taskStorage := &task.MongoStorage{
			BaseStorage: st,
			Codec:       task.Codec{},
		}

		return taskStorage, nil
	default:
		return nil, errx.UnsupportedKind(config.Kind)
	}
}

func BuildResources(ctx context.Context, broker *event.Broker, logger *slog.Logger, config ResourcesConfig) (resource.Storage, error) {
	storage, err := BuildResourcesStorage(ctx, config.Storage)
	if err != nil {
		return nil, errx.FailedToBuild("resources_storage", err)
	}

	for _, observerConfig := range config.Observers {
		observer, err := BuildResourcesObserver(ctx, broker, logger, observerConfig)
		if err != nil {
			return nil, errx.FailedToBuild("resources_observer", err)
		}

		storage = resources.Observe(storage, observer)
	}

	return storage, nil
}

func BuildResourcesStorage(ctx context.Context, config ResourcesStorageConfig) (resource.Storage, error) {
	switch config.Kind {
	case "redis":
		redisClient, err := buildRedis(ctx, config.URI)
		if err != nil {
			return nil, errx.FailedToBuild("redis", err)
		}

		keyPrefix, err := getOption[string](config.Options, "key_prefix")
		if err != nil {
			return nil, errx.FailedToBuild("redis", err)
		}

		return redisResources.New(redisClient, codec.JSON[resource.Resource](), keyPrefix), nil
	default:
		return nil, errx.UnsupportedKind(config.Kind)
	}
}

func BuildResourcesObserver(ctx context.Context, broker *event.Broker, logger *slog.Logger, config ResourcesObserverConfig) (resources.Observer, error) {
	switch config.Kind {
	case "notifier":
		var notifierConfig struct {
			Topic string `json:"topic" yaml:"topic"`
		}
		if err := config.Options.Bind(&notifierConfig); err != nil {
			return nil, errx.FailedToBind(err)
		}

		notifier := resources.Notifier{
			Logger: logger.With(slog.String("component", "resources_observer.notifier")),
			Broker: broker,
			Topic:  notifierConfig.Topic,
		}

		return notifier, nil
	default:
		return nil, errx.UnsupportedKind(config.Kind)
	}
}

func BuildEvents(ctx context.Context, logger *slog.Logger, config EventsConfig) (*event.Broker, error) {
	return BuildEventBroker(ctx, logger, config.Broker)
}

func BuildEventBroker(ctx context.Context, logger *slog.Logger, cfg EventsBrokerConfig) (*event.Broker, error) {
	switch cfg.Kind {
	case "nats":
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
			cfg.URI,
			subscriberConfig,
			publishedConfig,
			logger.With(
				slog.String("component", "broker_agent"),
				slog.String("component_type", "nats"),
			),
			//extractLogger(ctx, slog.Default()),
		)
		if err != nil {
			return nil, fmt.Errorf("failed to connect to nats: %w", err)
		}

		b := watermill.Broker[event.Event]{
			Agent:     agent,
			ItemCodec: codec.JSON[event.Event](),
			Logger: logger.With(
				slog.String("component", "broker"),
			),
			//Logger:                    extractLogger(ctx, slog.Default()),
			ConsumerChannelBufferSize: 1024,
		}

		return event.NewBroker(b), nil
	default:
		return nil, errx.UnsupportedKind(cfg.Kind)
	}
}

func BuildSpecifications(ctx context.Context, cfg SpecificationsConfig) (*specification.Manager, error) {
	switch cfg.Storage.Kind {
	case "postgres":
		pool, err := buildPostgres(ctx, cfg.Storage.URI)
		if err != nil {
			return nil, errx.FailedToBuild("postgres", err)
		}

		specificationsTable, err := getOption[string](cfg.Storage.Options, "specifications_table")
		if err != nil {
			return nil, err
		}

		typesTable, err := getOption[string](cfg.Storage.Options, "types_table")
		if err != nil {
			return nil, err
		}

		return specification.NewManager(
			&specification.PostgresBinaryStorage{
				DB:    pool,
				Table: specificationsTable,
			},
			&specification.PostgresBinaryStorage{
				DB:    pool,
				Table: typesTable,
			},
		), nil
	default:
		return nil, errx.UnsupportedKind(cfg.Storage.Kind)
	}
}

func BuildHttpApiDeployment(ctx context.Context, sys *system.System, specificationManager *specification.Manager, logger *slog.Logger, cfg HttpApiServiceConfig) (Deployment[*HttpApiService], error) {
	core, err := BuildServiceCore(ctx, "http-api", logger, cfg.ServiceConfig)
	if err != nil {
		return Deployment[*HttpApiService]{}, errx.FailedToBuild("core", err)
	}

	srvc := &HttpApiService{
		sys:            sys,
		specifications: specificationManager,
		port:           cfg.Port,
		loggerEnabled:  cfg.LoggerEnabled,
		Core:           core,
	}

	middlewares := buildMiddlewares(sys, cfg.ServiceConfig)

	return Deployment[*HttpApiService]{
		Service:     srvc,
		Middlewares: middlewares,
	}, nil
}

func BuildSchedulerDeployment(ctx context.Context, sys *system.System, logger *slog.Logger, cfg SchedulerServiceConfig) (Deployment[*dependencies.Service], error) {
	core, err := BuildServiceCore(ctx, "scheduler", logger, cfg.ServiceConfig)
	if err != nil {
		return Deployment[*dependencies.Service]{}, errx.FailedToBuild("core", err)
	}

	dependencyManager, err := buildDependencyManager(ctx, logger, cfg.Dependencies)
	if err != nil {
		return Deployment[*dependencies.Service]{}, errx.FailedToBuild("dependency_manager", err)
	}

	taskToGroup, err := buildTaskToGroup(ctx, cfg.TaskToGroup)
	if err != nil {
		return Deployment[*dependencies.Service]{}, errx.FailedToBuild("task2group", err)
	}

	resolvers, err := buildResolvers(ctx, sys, logger, cfg.Resolvers)
	if err != nil {
		return Deployment[*dependencies.Service]{}, errx.FailedToBuild("resolvers", err)
	}

	mngr := &manager2.Manager{
		System:       sys,
		Dependencies: dependencyManager,
		TaskToGroup:  taskToGroup,
		Resolvers:    resolvers,
		Logger:       logger.With(slog.String("component", "dependency_manager")),
		//Logger:       extractLogger(ctx, slog.Default()),
	}

	srvc := &dependencies.Service{
		System:  sys,
		Manager: mngr,
		Core:    core,
	}

	return Deployment[*dependencies.Service]{
		Service:     srvc,
		Middlewares: buildMiddlewares(sys, cfg.ServiceConfig),
	}, nil
}

func buildResolvers(ctx context.Context, system system.AbstractSystem, logger *slog.Logger, configs []SchedulerResolverConfig) (map[string]manager2.Resolver, error) {
	result := make(map[string]manager2.Resolver, len(configs))
	for _, config := range configs {
		switch config.Kind {
		case "resource":
			var resourceResolverConfig SchedulerResourceResolverConfig
			if err := config.Data.Bind(&resourceResolverConfig); err != nil {
				return nil, errx.FailedToBind(err)
			}

			resolver, err := buildResourceResolver(ctx, system, logger, resourceResolverConfig)
			if err != nil {
				return nil, errx.FailedToBuild("resource_resolver", err)
			}

			result["resource"] = resolver
		default:
			return nil, errx.UnsupportedKind(config.Kind)
		}
	}

	return result, nil
}

func buildResourceResolver(ctx context.Context, system system.AbstractSystem, logger *slog.Logger, cfg SchedulerResourceResolverConfig) (*resourceResolver2.Resolver, error) {
	storage, err := buildResourceResolverStorage(ctx, cfg.Storage)
	if err != nil {
		return nil, errx.FailedToBuild("resource_resolver_storage", err)
	}

	resolver := &resourceResolver2.Resolver{
		System:       system,
		Storage:      storage,
		PollLimit:    cfg.Poller.Limit,
		PollInterval: cfg.Poller.Interval,
		Logger: logger.With(
			slog.String("component", "dependency_resolver"),
			slog.String("component_type", "resource"),
		),
		//Logger:       extractLogger(ctx, slog.Default()),
	}

	return resolver, nil
}

func buildResourceResolverStorage(ctx context.Context, cfg SchedulerResourceResolverStorageConfig) (resourceResolver2.Storage, error) {
	switch cfg.Kind {
	case "mongo":
		return buildMongoResourceResolverStorage(ctx, cfg)
	default:
		return nil, errx.UnsupportedKind(cfg.Kind)
	}
}

func buildMongoResourceResolverStorage(ctx context.Context, cfg SchedulerResourceResolverStorageConfig) (*resourceResolver2.MongoStorage, error) {
	conn, err := buildMongo(ctx, cfg.URI)
	if err != nil {
		return nil, errx.FailedToBuild("mongo", err)
	}

	var mongoResourseResolverStorageConfig struct {
		DB          string        `yaml:"db,omitempty" json:"db,omitempty"`
		Collection  string        `yaml:"collection,omitempty" json:"collection,omitempty"`
		PollTimeout time.Duration `yaml:"poll_timeout,omitempty" json:"poll_timeout,omitempty"`
	}
	dc := DynamicConfig(cfg.Options)
	if err := dc.Bind(&mongoResourseResolverStorageConfig); err != nil {
		return nil, errx.FailedToBind(err)
	}

	return &resourceResolver2.MongoStorage{
		Collection: conn.
			Database(mongoResourseResolverStorageConfig.DB).
			Collection(mongoResourseResolverStorageConfig.Collection),
		PollTimeout: mongoResourseResolverStorageConfig.PollTimeout,
	}, nil
}

func buildTaskToGroup(ctx context.Context, cfg SchedulerTaskToGroupConfig) (manager2.TaskToGroup, error) {
	switch cfg.Kind {
	case "redis":
		client, err := buildRedis(ctx, cfg.URI)
		if err != nil {
			return nil, errx.FailedToBuild("redis", err)
		}
		return &task2group.RedisStorage{Client: client}, nil
	default:
		return nil, errx.UnsupportedKind(cfg.Kind)
	}
}

func buildDependencyManager(ctx context.Context, logger *slog.Logger, cfg SchedulerDependenciesConfig) (dependency.Manager, error) {
	return buildBatchedPostgresDependencies(ctx, logger, cfg)
}

func buildBatchedPostgresDependencies(ctx context.Context, logger *slog.Logger, cfg SchedulerDependenciesConfig) (*batched2.Manager, error) {
	pg, err := buildPostgres(ctx, cfg.Postgres.URI)
	if err != nil {
		return nil, fmt.Errorf("failed to build: %w", err)
	}

	mng := &batched2.Manager{
		Client: pg,
		Config: batched2.Config{
			PollingInterval:  cfg.Poller.Interval,
			PollingBatchSize: cfg.Poller.BatchSize,
		},
		Logger: logger.With(
			slog.String("component", "dependencies"),
			slog.String("component_type", "postgres:batched"),
		),
		//Logger: extractLogger(ctx, slog.Default()),
	}

	return mng, nil
}

func BuildStatusDeployment(ctx context.Context, sys *system.System, logger *slog.Logger, cfg StatusServiceConfig) (Deployment[*status.Service], error) {
	core, err := BuildServiceCore(ctx, "status", logger, cfg.ServiceConfig)
	if err != nil {
		return Deployment[*status.Service]{}, errx.FailedToBuild("core", err)
	}

	middlewares := buildMiddlewares(sys, cfg.ServiceConfig)

	srvc := &status.Service{
		System:      sys,
		ResultCodec: codec.JSON[executor.Result](),
		Core:        core,
	}

	return Deployment[*status.Service]{
		Service:     srvc,
		Middlewares: middlewares,
	}, nil
}

func BuildProcessorDeployment(ctx context.Context, sys *system.System, exe executor.Executor, logger *slog.Logger, cfg ProcessorServiceConfig) (Deployment[*executor.Service], error) {
	core, err := BuildServiceCore(ctx, "processor", logger, cfg.ServiceConfig)
	if err != nil {
		return Deployment[*executor.Service]{}, errx.FailedToBuild("core", err)
	}

	middlewares := buildMiddlewares(sys, cfg.ServiceConfig)

	srvc := &executor.Service{
		System:      sys,
		ResultCodec: codec.JSON[executor.Result](),
		Executor:    exe,
		Core:        core,
	}

	return Deployment[*executor.Service]{
		Service:     srvc,
		Middlewares: middlewares,
	}, nil
}

func BuildDiscoveryDeployment(ctx context.Context, sys *system.System, logger *slog.Logger, cfg DiscoveryServiceConfig) (Deployment[*discovery.Poller], error) {
	core, err := BuildServiceCore(ctx, "discovery", logger, cfg.ServiceConfig)
	if err != nil {
		return Deployment[*discovery.Poller]{}, errx.FailedToBuild("core", err)
	}

	mws := buildMiddlewares(sys, cfg.ServiceConfig)

	hub, err := buildDiscoveryHub(ctx, cfg.Hub)
	if err != nil {
		return Deployment[*discovery.Poller]{}, errx.FailedToBuild("hub", err)
	}

	srvc := &discovery.Poller{
		Hub:           hub,
		Events:        sys.Events(),
		RequestCodec:  codec.JSON[discovery.Request](),
		ResponseCodec: codec.JSON[discovery.Response](),
		Interval:      cfg.PollingInterval,
		Core:          core,
	}

	return Deployment[*discovery.Poller]{
		Service:     srvc,
		Middlewares: mws,
	}, nil
}

func buildDiscoveryHub(ctx context.Context, cfg DiscoveryHubConfig) (discovery.Hub, error) {
	switch cfg.Kind {
	case "consul":
		return buildConsulHub(ctx, cfg.Options)
	default:
		return nil, errx.UnsupportedKind(cfg.Kind)
	}
}

func buildConsulHub(ctx context.Context, cfg DynamicConfig) (*consul.Hub, error) {
	var consulConfig struct {
		Addr string `yaml:"addr,omitempty" json:"addr,omitempty"`
	}
	if err := cfg.Bind(&consulConfig); err != nil {
		return nil, errx.FailedToBind(err)
	}

	client, err := capi.NewClient(&capi.Config{
		Address: consulConfig.Addr,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create a consul client: %w", err)
	}

	return &consul.Hub{Consul: client}, nil
}

func BuildServiceCore(ctx context.Context, defaultName string, logger *slog.Logger, config ServiceConfig) (service.Core, error) {
	if config.Name == "" {
		config.Name = defaultName
	}

	if config.ID == "" {
		config.ID = uid.Generate()
	}

	core := service.NewCore(
		config.Name,
		config.ID,
		logger.With(
			slog.String("deployment_id", config.ID),
		),
		//newLogger(os.Stdout).With("service", data.Info.Name),
	)

	return core, nil
}

func BuildServiceLogger(
	handler slog.Handler,
	service string,
) *slog.Logger {

	logger := slog.New(handler).With("service", service)

	return logger
}

func BuildPrettySlogHandler(w io.Writer, level slog.Level) slog.Handler {
	coloredHandler := tint.NewHandler(w, &tint.Options{
		Level:      level,
		TimeFormat: time.DateTime,
	})

	prefixedHandler := prefixed.NewHandler(coloredHandler,
		&prefixed.HandlerOptions{
			PrefixKeys: []string{"time", "service"},
		})

	return prefixedHandler
}

func buildMiddlewares(sys *system.System, cfg ServiceConfig) []service.Middleware {
	var middlewares []service.Middleware

	if cfg.Discovery.Enabled {
		middlewares = append(middlewares, discovery.WithStaticInfo[service.Service](
			cfg.Discovery.Info,
			sys.Events(),
			codec.JSON[discovery.Request](),
			codec.JSON[discovery.Response](),
		))
	}

	middlewares = append(middlewares, loggingMiddleware{})

	return middlewares
}

func buildMongo(ctx context.Context, uri string) (*mongo.Client, error) {
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
	if err != nil {
		return nil, fmt.Errorf("failed to connect: %w", err)
	}

	return client, nil
}

func buildRedis(ctx context.Context, uri string) (redis.UniversalClient, error) {
	client := redis.NewUniversalClient(&redis.UniversalOptions{
		Addrs: []string{uri},
	})

	return client, nil
}

func buildPostgres(ctx context.Context, uri string) (*pgxpool.Pool, error) {
	pool, err := pgxpool.New(context.Background(), uri)
	if err != nil {
		return nil, fmt.Errorf("failed to create a pool: %w", err)
	}

	return pool, nil
}

func getOption[T any](options map[string]any, key string) (T, error) {
	var zero T

	raw, ok := options[key]
	if !ok {
		return zero, fmt.Errorf("no option %s found", key)
	}

	val, ok := raw.(T)
	if !ok {
		return zero, fmt.Errorf("option %s type mismatch", key)
	}

	return val, nil
}

package builder

import (
	"context"
	"fmt"
	"github.com/ischenkx/kantoku/pkg/common/dependency"
	"github.com/ischenkx/kantoku/pkg/common/service"
	"github.com/ischenkx/kantoku/pkg/core/services/scheduler/dependencies/simple"
	"github.com/ischenkx/kantoku/pkg/core/services/scheduler/dependencies/simple/manager"
	resourceResolver "github.com/ischenkx/kantoku/pkg/core/services/scheduler/dependencies/simple/manager/resolvers/resource_resolver"
	"github.com/ischenkx/kantoku/pkg/core/services/scheduler/dependencies/simple/manager/task2group"
	"github.com/ischenkx/kantoku/pkg/core/services/scheduler/dummy"
	"github.com/ischenkx/kantoku/pkg/core/system"
	"github.com/ischenkx/kantoku/pkg/lib/gateway/cli/config"
	"github.com/ischenkx/kantoku/pkg/lib/gateway/cli/errx"
	"github.com/ischenkx/kantoku/pkg/lib/impl/data/dependency/postgres/batched"
	"log/slog"
	"time"
)

func (builder *Builder) BuildScheduler(ctx context.Context, sys system.AbstractSystem, cfg config.DynamicConfig) (srvc service.Service, mws []service.Middleware, err error) {
	core, err := builder.BuildServiceCore(ctx, "scheduler", cfg)
	if err != nil {
		return nil, nil, errx.FailedToBuild("core", err)
	}

	ctx = withLogger(ctx, core.Logger())

	switch cfg.Kind() {
	case "dependencies":
		mngr, err := builder.buildDependencyServiceManager(ctx, sys, cfg)
		if err != nil {
			return nil, nil, errx.FailedToBuild("dependencies service manager", err)
		}
		srvc = &simple.Service{
			System:  sys,
			Manager: mngr,
			Core:    core,
		}
	case "dummy":
		srvc = &dummy.Service{
			System: sys,
			Core:   core,
		}
	default:
		return nil, nil, fmt.Errorf("unsupported scheduler kind: %s", cfg.Kind())
	}

	middlewares, err := builder.BuildMiddlewares(ctx, srvc, sys, cfg)
	if err != nil {
		return nil, nil, errx.FailedToBuild("middlewares", err)
	}

	return srvc, middlewares, nil
}

func (builder *Builder) buildDependencyServiceManager(ctx context.Context, system system.AbstractSystem, cfg config.DynamicConfig) (*manager.Manager, error) {
	var dependenciesImplConfig struct {
		Task2Group   config.DynamicConfig
		Dependencies config.DynamicConfig
		Resolvers    config.DynamicConfig
	}
	if err := cfg.Bind(&dependenciesImplConfig); err != nil {
		return nil, errx.FailedToBind(err)
	}

	dependencyManager, err := builder.buildDependencyManager(ctx, dependenciesImplConfig.Dependencies)
	if err != nil {
		return nil, errx.FailedToBuild("dependency manager", err)
	}

	taskToGroup, err := builder.buildTaskToGroup(ctx, dependenciesImplConfig.Task2Group)
	if err != nil {
		return nil, errx.FailedToBuild("task2group", err)
	}

	resolvers, err := builder.buildResolvers(ctx, system, dependenciesImplConfig.Resolvers)
	if err != nil {
		return nil, errx.FailedToBuild("resolvers", err)
	}

	mng := &manager.Manager{
		System:       system,
		Dependencies: dependencyManager,
		TaskToGroup:  taskToGroup,
		Resolvers:    resolvers,
		Logger:       extractLogger(ctx, slog.Default()),
	}

	return mng, nil
}

func (builder *Builder) buildResolvers(ctx context.Context, system system.AbstractSystem, cfg config.DynamicConfig) (map[string]manager.Resolver, error) {
	var resolvers struct {
		Resource config.DynamicConfig
	}
	if err := cfg.Bind(&resolvers); err != nil {
		return nil, errx.FailedToBind(err)
	}

	_resourceResolver, err := builder.buildResourceResolver(ctx, system, resolvers.Resource)
	if err != nil {
		return nil, errx.FailedToBuild("resource resolver", err)
	}

	mapping := map[string]manager.Resolver{
		"resource": _resourceResolver,
	}

	return mapping, nil
}

func (builder *Builder) buildResourceResolver(ctx context.Context, system system.AbstractSystem, cfg config.DynamicConfig) (*resourceResolver.Resolver, error) {
	var resolverConfig struct {
		Storage config.DynamicConfig
		Poller  struct {
			Limit    int
			Interval time.Duration
		}
	}
	if err := cfg.Bind(&resolverConfig); err != nil {
		return nil, errx.FailedToBind(err)
	}

	storage, err := builder.buildResourceResolverStorage(ctx, resolverConfig.Storage)
	if err != nil {
		return nil, errx.FailedToBuild("resource resolver storage", err)
	}

	resolver := &resourceResolver.Resolver{
		System:       system,
		Storage:      storage,
		PollLimit:    resolverConfig.Poller.Limit,
		PollInterval: resolverConfig.Poller.Interval,
		Logger:       extractLogger(ctx, slog.Default()),
	}

	return resolver, nil
}

func (builder *Builder) buildResourceResolverStorage(ctx context.Context, cfg config.DynamicConfig) (resourceResolver.Storage, error) {
	switch cfg.Kind() {
	case "mongo":
		return builder.buildMongoResourceResolverStorage(ctx, cfg)
	default:
		return nil, errx.UnsupportedKind(cfg.Kind())
	}
}

func (builder *Builder) buildMongoResourceResolverStorage(ctx context.Context, cfg config.DynamicConfig) (*resourceResolver.MongoStorage, error) {
	var storageConfig struct {
		Conn        config.DynamicConfig
		PollTimeout time.Duration
	}
	if err := cfg.Bind(&storageConfig); err != nil {
		return nil, errx.FailedToBind(err)
	}

	mongoInfo, err := builder.BuildMongo(ctx, storageConfig.Conn)
	if err != nil {
		return nil, errx.FailedToBuild("mongo", err)
	}

	return &resourceResolver.MongoStorage{Collection: mongoInfo.GetCollection(), PollTimeout: storageConfig.PollTimeout}, nil
}

func (builder *Builder) buildTaskToGroup(ctx context.Context, cfg config.DynamicConfig) (manager.TaskToGroup, error) {
	switch cfg.Kind() {
	case "redis":
		client, err := builder.BuildRedis(ctx, cfg)
		if err != nil {
			return nil, errx.FailedToBuild("redis", err)
		}
		return &task2group.RedisStorage{Client: client}, nil
	default:
		return nil, errx.UnsupportedKind(cfg.Kind())
	}
}

func (builder *Builder) buildDependencyManager(ctx context.Context, cfg config.DynamicConfig) (dependency.Manager, error) {
	switch cfg.Kind() {
	case "postgres:batched":
		return builder.buildBatchedPostgresDependencies(ctx, cfg)
	default:
		return nil, errx.UnsupportedKind(cfg.Kind())
	}
}

func (builder *Builder) buildBatchedPostgresDependencies(ctx context.Context, cfg config.DynamicConfig) (*batched.Manager, error) {
	var managerConfig struct {
		Poller struct {
			Interval  time.Duration
			BatchSize int `mapstructure:"batch_size"`
		}
		Postgres config.DynamicConfig
	}
	if err := cfg.Bind(&managerConfig); err != nil {
		return nil, errx.FailedToBind(err)
	}

	pg, err := builder.BuildPostgres(ctx, managerConfig.Postgres)
	if err != nil {
		return nil, fmt.Errorf("failed to build: %w", err)
	}

	mng := &batched.Manager{
		Client: pg,
		Config: batched.Config{
			PollingInterval:  managerConfig.Poller.Interval,
			PollingBatchSize: managerConfig.Poller.BatchSize,
		},
		Logger: extractLogger(ctx, slog.Default()),
	}

	return mng, nil
}

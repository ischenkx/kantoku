package resourceResolverV2

import (
	"context"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/ischenkx/kantoku/pkg/common/data"
	"github.com/ischenkx/kantoku/pkg/common/transport/broker"
	"github.com/ischenkx/kantoku/pkg/common/transport/queue"
	"github.com/ischenkx/kantoku/pkg/core/event"
	"github.com/ischenkx/kantoku/pkg/core/resource"
	"github.com/ischenkx/kantoku/pkg/core/services/scheduler/dependencies/simple/manager"
	"github.com/ischenkx/kantoku/pkg/core/system"
	"log/slog"
)

type Resolver struct {
	System              system.AbstractSystem
	ReadyResourcesTopic string
	Logger              *slog.Logger
	Bindings            BindingStorage
}

func (resolver *Resolver) Bind(ctx context.Context, depId string, data any) (manager.BindingResult, error) {
	resourceId, ok := data.(string)
	if !ok {
		return manager.BindingResult{}, fmt.Errorf("unexpected data: %s", data)
	}

	if err := resolver.Bindings.Save(ctx, depId, resourceId); err != nil {
		return manager.BindingResult{}, fmt.Errorf("failed to bind resource and dependency: %w", err)
	}

	resolver.Logger.Info("binding saved",
		slog.String("dep_id", depId),
		slog.String("res_id", resourceId))

	res, err := resolver.System.Resources().Load(ctx, resourceId)
	if err != nil {
		return manager.BindingResult{}, fmt.Errorf("failed to load a resource: %w", err)
	}

	if res[0].Status == resource.Ready {
		fmt.Println("READY RESOURCE:", res[0].ID)
		return manager.BindingResult{Disabled: true}, nil
	}

	return manager.BindingResult{Disabled: false}, nil
}

func (resolver *Resolver) Ready(ctx context.Context) (<-chan string, error) {
	events, err := resolver.System.Events().Consume(ctx, broker.TopicsInfo{
		Group:  "resource_resolver",
		Topics: []string{resolver.ReadyResourcesTopic},
	})
	if err != nil {
		return nil, fmt.Errorf("failed to consume resource events: %w", err)
	}

	depIds := make(chan string, 1024)

	go resolver.collectResolvedDependencies(ctx, events, depIds)

	return depIds, nil
}

func (resolver *Resolver) collectResolvedDependencies(ctx context.Context, events <-chan queue.Message[event.Event], ids chan<- string) {
	queue.Processor[event.Event]{
		Handler: func(ctx context.Context, ev event.Event) error {
			resourceId := string(ev.Data)
			resolver.Logger.Info("(resource resolver) received a resource",
				slog.String("id", resourceId))
			dependencyId, err := resolver.Bindings.Load(ctx, resourceId)
			if err != nil {
				if errors.Is(err, data.NotFoundErr) {
					resolver.Logger.Info("(resource resolver) resource not found",
						slog.String("id", resourceId))
					return nil
				}
				return fmt.Errorf("failed to load a dependency id: %w", err)
			}

			id := uuid.New().String()

			fmt.Println("here we go!", id)

			select {
			case <-ctx.Done():
				return fmt.Errorf("context done")
			case ids <- dependencyId:
			}

			fmt.Println("DONE!", id)

			return nil
		},
		ErrorHandler: func(ctx context.Context, ev event.Event, err error) {
			resourceId := string(ev.Data)

			resolver.Logger.Error("failed to process a ready resource",
				slog.String("id", resourceId),
				slog.String("error", err.Error()))
		},
	}.Process(ctx, events)
}

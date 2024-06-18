package resourceResolver

import (
	"context"
	"fmt"
	"github.com/ischenkx/kantoku/pkg/core/resource"
	"github.com/ischenkx/kantoku/pkg/core/system"
	"github.com/samber/lo"
	"log/slog"
	"time"
)

type Resolver struct {
	System       system.AbstractSystem
	Storage      Storage
	PollLimit    int
	PollInterval time.Duration
	Logger       *slog.Logger
}

func (resolver *Resolver) Bind(ctx context.Context, depId string, data any) error {
	resourceId, ok := data.(string)
	if !ok {
		return fmt.Errorf("unexpected data: %s", data)
	}

	if err := resolver.Storage.Save(ctx, depId, resourceId); err != nil {
		return fmt.Errorf("failed to bind resource and dependency: %w", err)
	}

	return nil
}

func (resolver *Resolver) Ready(ctx context.Context) (<-chan string, error) {
	depIds := make(chan string, 1024)

	go resolver.collectResolvedDependencies(ctx, depIds)

	return depIds, nil
}

func (resolver *Resolver) collectResolvedDependencies(ctx context.Context, ids chan<- string) {

	pollLimit := resolver.PollLimit
	if pollLimit <= 0 {
		pollLimit = 1024
	}

	pollInterval := resolver.PollInterval
	if pollInterval <= 0 {
		pollInterval = time.Second * 5
	}

	ticker := time.NewTicker(pollInterval)
	defer ticker.Stop()
	resolver.Logger.Info("collecting resolved dependencies",
		slog.Duration("interval", pollInterval))

poller:
	for {
		select {
		case <-ctx.Done():
			break poller

		case <-ticker.C:
			bindings, err := resolver.Storage.Poll(ctx, pollLimit)
			if err != nil {
				resolver.Logger.Error("failed to poll bindings",
					slog.String("error", err.Error()))
				continue
			}

			resource2dependencies := lo.GroupBy(bindings, func(binding Binding) string {
				return binding.ResourceId
			})

			resourceIds := lo.Keys(resource2dependencies)
			resources, err := resolver.System.Resources().Load(ctx, resourceIds...)
			if err != nil {
				resolver.Logger.Error("failed to load resources",
					slog.String("error", err.Error()))
			}

			var resolvedIds []string

			for _, res := range resources {
				if res.Status != resource.Ready {
					continue
				}

				resourceBindings := resource2dependencies[res.ID]
				for _, binding := range resourceBindings {
					select {
					case <-ctx.Done():
						break poller
					case ids <- binding.DependencyId:
						resolvedIds = append(resolvedIds, binding.DependencyId)
					}
				}
			}
			if err := resolver.Storage.Resolve(ctx, resolvedIds...); err != nil {
				resolver.Logger.Error("failed to resolve resources",
					slog.String("error", err.Error()))
			}
		}
	}
}

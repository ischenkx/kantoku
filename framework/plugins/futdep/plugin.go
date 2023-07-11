package futdep

import (
	"context"
	"kantoku/common/data/pool"
	"kantoku/framework/future"
	"kantoku/framework/infra/demon"
	"kantoku/framework/plugins/depot"
	"kantoku/framework/utils/demons"
	"kantoku/kernel"
	"log"
)

type Plugin struct {
	manager *Manager
	futures *future.Manager
}

func NewPlugin(manager *Manager, futures *future.Manager) *Plugin {
	return &Plugin{
		manager: manager,
		futures: futures,
	}
}

type PluginData struct {
	Futures []future.ID
}

func (p *Plugin) BeforeScheduled(ctx *kernel.Context) error {
	data := ctx.Data().GetOrSet("futdep", func() any { return &PluginData{} }).(*PluginData)

	for _, subtask := range data.Futures {
		dependency, err := p.manager.Make(ctx, subtask)
		if err != nil {
			return err
		}

		if err := depot.Dependency(dependency)(ctx); err != nil {
			return err
		}
	}

	return nil
}

func (p *Plugin) Demons(ctx context.Context) []demon.Demon {
	return demons.Multi{
		demons.TryProvider(p.manager.deps),
		demons.TryProvider(p.manager.fut2dep),
		demons.Functional("FUTDEP_PROCESSOR", p.process),
	}.Demons(ctx)
}

func (p *Plugin) process(ctx context.Context) error {
	return pool.ReadAutoCommit[future.ID](ctx, p.futures.Resolutions(), func(ctx context.Context, id future.ID) error {
		if err := p.manager.ResolveFuture(ctx, id); err != nil {
			log.Println("failed to resolve a future:", err)
		}
		return nil
	})
}

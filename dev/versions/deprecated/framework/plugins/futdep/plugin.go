package futdep

import (
	"context"
	future "kantoku/common/data/future"
	"kantoku/common/data/pool"
	demons2 "kantoku/common/util/demons"
	"kantoku/framework/infra"
	"kantoku/framework/job"
	"kantoku/framework/plugins/depot"
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

func (p *Plugin) BeforeScheduled(ctx *job.Context) error {
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

func (p *Plugin) Demons() []infra.Demon {
	return demons2.Multi{
		demons2.TryProvider(p.manager.deps),
		demons2.TryProvider(p.manager.fut2dep),
		demons2.Functional("FUTDEP_PROCESSOR", p.process),
	}.Demons()
}

func (p *Plugin) process(ctx context.Context) error {
	return pool.AutoCommit[future.ID](ctx, p.futures.Resolutions(), func(ctx context.Context, id future.ID) error {
		if err := p.manager.ResolveFuture(ctx, id); err != nil {
			log.Println("failed to resolve a future:", err)
		}
		return nil
	})
}

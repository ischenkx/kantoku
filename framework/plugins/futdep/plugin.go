package futdep

import (
	"kantoku/framework/future"
	"kantoku/framework/plugins/depot"
	"kantoku/kernel"
)

type Plugin struct {
	manager *Manager
}

func NewPlugin(manager *Manager) *Plugin {
	return &Plugin{
		manager: manager,
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

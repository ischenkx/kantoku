package delay

import (
	"kantoku"
	"kantoku/framework/depot"
	"time"
)

type Plugin struct {
	manager *Manager
}

func NewPlugin(manager *Manager) *Plugin {
	return &Plugin{manager: manager}
}

type PluginData struct {
	When *time.Time
}

func (p *Plugin) Initialize(kantoku *kantoku.Kantoku) {}

func (p *Plugin) BeforeInitialized(ctx *kantoku.Context) error {
	return nil
}

func (p *Plugin) AfterInitialized(ctx *kantoku.Context) {}

func (p *Plugin) BeforeScheduled(ctx *kantoku.Context) error {
	data := ctx.Data().GetOrSet("delay", func() any { return &PluginData{} }).(*PluginData)

	if data.When != nil {
		dep, err := p.manager.MakeDependency(ctx, *data.When)
		if err != nil {
			return err
		}

		depot.Dependency(dep)(ctx)
	}

	return nil
}

func (p *Plugin) AfterScheduled(ctx *kantoku.Context) {}

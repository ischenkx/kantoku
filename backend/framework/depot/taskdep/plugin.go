package taskdep

import (
	"kantoku"
	"kantoku/backend/framework/depot"
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
	Subtasks []string
}

func (p *Plugin) Initialize(kantoku *kantoku.Kantoku) {}

func (p *Plugin) BeforeInitialized(ctx *kantoku.Context) error {
	return nil
}

func (p *Plugin) AfterInitialized(ctx *kantoku.Context) {}

func (p *Plugin) BeforeScheduled(ctx *kantoku.Context) error {
	data := ctx.Data().GetOrSet("taskdep", func() any { return &PluginData{} }).(*PluginData)

	for _, subtask := range data.Subtasks {
		dependency, err := p.manager.SubtaskDependency(ctx, subtask)
		if err != nil {
			return err
		}

		if err := depot.Dependency(dependency)(ctx); err != nil {
			return err
		}
	}

	return nil
}

func (p *Plugin) AfterScheduled(ctx *kantoku.Context) {}

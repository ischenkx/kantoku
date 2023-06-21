package taskdep

import (
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
	Subtasks []string
}

func (p *Plugin) BeforeScheduled(ctx *kernel.Context) error {
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

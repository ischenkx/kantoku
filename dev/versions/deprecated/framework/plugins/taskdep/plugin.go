package taskdep

import (
	"context"
	"kantoku/backend/executor"
	demons2 "kantoku/common/util/demons"
	"kantoku/framework/infra"
	job2 "kantoku/framework/job"
	"kantoku/framework/plugins/depot"
	"log"
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

func (p *Plugin) BeforeScheduled(ctx *job2.Context) error {
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

func (p *Plugin) Demons() []infra.Demon {
	return demons2.Multi{
		demons2.TryProvider(p.manager.deps),
		demons2.TryProvider(p.manager.task2dep),
	}.Demons()
}

func (p *Plugin) ExecutorPlugins() []executor.Plugin {
	return []executor.Plugin{
		ExecutorPlugin{manager: p.manager},
	}
}

type ExecutorPlugin struct {
	manager *Manager
}

func (plugin ExecutorPlugin) SavedTaskResult(ctx context.Context, id string, result job2.Result) {
	if err := plugin.manager.ResolveTask(ctx, id); err != nil {
		log.Println("failed to resolve task dependency:", err)
	}
}

func (plugin ExecutorPlugin) FailedToSaveTaskResult(ctx context.Context, id string, result job2.Result) {
	log.Println("failed to save task result:", id)
}

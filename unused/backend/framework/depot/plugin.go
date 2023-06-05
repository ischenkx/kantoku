package depot

import (
	"context"
	"github.com/samber/lo"
	"kantoku"
	"kantoku/unused/backend/framework/depot/deps"
)

type Plugin struct {
	depot *Depot
}

type PluginData struct {
	Dependencies []string
}

func NewPlugin(depot *Depot) *Plugin {
	return &Plugin{depot: depot}
}

func (p *Plugin) BeforeInitialized(ctx *kantoku.Context) error {
	return nil
}

func (p *Plugin) AfterInitialized(ctx *kantoku.Context) {

}

func (p *Plugin) BeforeScheduled(ctx *kantoku.Context) error {
	return nil
}

func (p *Plugin) AfterScheduled(ctx *kantoku.Context) {

}

type DependenciesEvaluator struct {
	depot *Depot
}

func (e DependenciesEvaluator) Evaluate(ctx context.Context, task string) (any, error) {
	groupID, err := e.depot.GroupTaskBimap().ByValue(ctx, task)
	if err != nil {
		return nil, err
	}

	group, err := e.depot.Deps().Group(ctx, groupID)
	if err != nil {
		return nil, err
	}

	dep2status := lo.Associate(group.Dependencies, func(dep deps.Dependency) (string, bool) {
		return dep.ID, dep.Resolved
	})

	return dep2status, nil
}

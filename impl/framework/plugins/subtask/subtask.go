package subtask

import (
	"context"
	"kantoku"
	"kantoku/common/deps"
	"kantoku/framework/plugins"
	plugins2 "kantoku/impl/framework/plugins"
)

type Subtask struct {
	deps deps.Deps
}

func New(deps deps.Deps) Subtask {
	return Subtask{deps: deps}
}

func (s Subtask) Id() string {
	return plugins2.Id()
}

func Output(task kantoku.Task) plugins.ArgumentCallback {
	return Callback{subtask: task}
}

func (s Subtask) DependencyFromTask(ctx context.Context, task kantoku.Task) (*deps.Dependency, error) {
	dep, err := s.deps.Make(ctx)
	if err != nil {
		return nil, err
	}
	// TODO: CRITICAL! check if task is already resolved
	return dep, nil
}

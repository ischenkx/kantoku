package taskdep

import (
	"context"
	"kantoku/common/data"
	"kantoku/common/data/kv"
	"kantoku/framework/plugins/depot/deps"
	"log"
)

type Manager struct {
	task2dep kv.Database[string, string]
	deps     deps.Deps
}

func NewManager(deps deps.Deps, task2dep kv.Database[string, string]) *Manager {
	return &Manager{
		task2dep: task2dep,
		deps:     deps,
	}
}

func (manager *Manager) SubtaskDependency(ctx context.Context, task string) (string, error) {
	dependencyID, err := manager.task2dep.Get(ctx, task)

	if err == nil {
		return dependencyID, nil
	} else if err != data.NotFoundErr {
		return "", err
	}

	dep, err := manager.deps.Make(ctx)
	if err != nil {
		return "", err
	}

	log.Println("Making a taskdep dep:", dep.ID, task)

	res, _, err := manager.task2dep.GetOrSet(ctx, task, dep.ID)
	return res, err
}

func (manager *Manager) ResolveTask(ctx context.Context, id string) error {
	dep, err := manager.task2dep.Get(ctx, id)
	if err != nil {
		return err
	}
	return manager.deps.Resolve(ctx, dep)
}

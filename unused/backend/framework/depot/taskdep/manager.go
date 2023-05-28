package taskdep

import (
	"context"
	"kantoku/common/data"
	"kantoku/common/data/kv"
	"kantoku/unused/backend/framework/depot/deps"
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

	return manager.task2dep.GetOrSet(ctx, task, dep.ID)
}

func (manager *Manager) ResolveTask(ctx context.Context, id string) error {
	log.Println("resolving:", id, manager.task2dep)
	dep, err := manager.task2dep.Get(ctx, id)
	if err != nil {
		return err
	}
	return manager.deps.Resolve(ctx, dep)
}

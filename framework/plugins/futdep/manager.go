package futdep

import (
	"context"
	"fmt"
	"kantoku/common/data"
	"kantoku/common/data/kv"
	"kantoku/framework/future"
	"kantoku/framework/plugins/depot/deps"
)

type Manager struct {
	fut2dep kv.Database[future.ID, string]
	deps    deps.Deps
}

func NewManager(deps deps.Deps, fut2dep kv.Database[future.ID, string]) *Manager {
	return &Manager{
		fut2dep: fut2dep,
		deps:    deps,
	}
}

func (manager *Manager) Make(ctx context.Context, id future.ID) (string, error) {
	dependencyID, err := manager.fut2dep.Get(ctx, id)

	if err == nil {
		return dependencyID, nil
	} else if err != data.NotFoundErr {
		return "", fmt.Errorf("failed to get data from fut2dep: %s", err)
	}

	dep, err := manager.deps.NewDependency(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to make a dependency: %s", err)
	}
	res, _, err := manager.fut2dep.GetOrSet(ctx, id, dep.ID)
	if err != nil {
		return "", fmt.Errorf("failed to retrieve a dependency by a future id: %s", err)
	}
	return res, err
}

func (manager *Manager) ResolveFuture(ctx context.Context, id future.ID) error {
	dep, err := manager.fut2dep.Get(ctx, id)
	if err != nil {
		return fmt.Errorf("failed to get a dependency id for the given future: %s", err)
	}
	return manager.deps.Resolve(ctx, dep)
}

package futdep

import (
	"context"
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
		return "", err
	}

	dep, err := manager.deps.Make(ctx)
	if err != nil {
		return "", err
	}
	res, _, err := manager.fut2dep.GetOrSet(ctx, id, dep.ID)
	return res, err
}

func (manager *Manager) ResolveFuture(ctx context.Context, id future.ID) error {
	dep, err := manager.fut2dep.Get(ctx, id)
	if err != nil {
		return err
	}
	return manager.deps.Resolve(ctx, dep)
}

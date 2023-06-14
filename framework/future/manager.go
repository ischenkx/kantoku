package future

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"kantoku/common/data/kv"
)

type Manager struct {
	futures   kv.Database[ID, Future]
	resources kv.Database[ID, Resource]
	runner    Runner
}

func NewManager(futures kv.Database[ID, Future], resources kv.Database[ID, Resource], runner Runner) *Manager {
	return &Manager{
		futures:   futures,
		resources: resources,
		runner:    runner,
	}
}

func (manager *Manager) Make(ctx context.Context, typ string, param []byte) (Future, error) {
	id := ID(uuid.New().String())
	future := Future{
		ID:    id,
		Type:  typ,
		Param: param,
	}
	err := manager.futures.Set(ctx, id, future)
	return future, err
}

func (manager *Manager) Resolve(ctx context.Context, id ID, resource Resource) error {
	future, err := manager.futures.Get(ctx, id)
	if err != nil {
		return err
	}
	resource, set, err := manager.resources.GetOrSet(ctx, id, resource)
	if err != nil {
		return err
	}
	if !set {
		return ErrAlreadyResolved
	}
	manager.runner.Run(ctx, Resolution{Future: future, Resource: resource})
	return nil
}

func (manager *Manager) Load(ctx context.Context, id ID) (Resolution, error) {
	resource, err := manager.resources.Get(ctx, id)
	if err != nil {
		return Resolution{}, fmt.Errorf("failed to load a resource: %s", err)
	}

	future, err := manager.futures.Get(ctx, id)
	if err != nil {
		return Resolution{}, fmt.Errorf("failed to load a future: %s", err)
	}

	return Resolution{Future: future, Resource: resource}, nil
}

package future

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"kantoku/common/data"
	"kantoku/common/data/kv"
	"kantoku/common/data/pool"
)

type Manager struct {
	futures     kv.Database[ID, Future]
	resources   kv.Database[ID, Resource]
	resolutions pool.Pool[ID]
}

func NewManager(futures kv.Database[ID, Future], resources kv.Database[ID, Resource], resolutions pool.Pool[ID]) *Manager {
	return &Manager{
		futures:     futures,
		resources:   resources,
		resolutions: resolutions,
	}
}

func (manager *Manager) Make(ctx context.Context, typ string, param []byte) (Future, error) {
	id := uuid.New().String()
	future := Future{
		ID:    id,
		Type:  typ,
		Param: param,
	}
	err := manager.futures.Set(ctx, id, future)
	return future, err
}

func (manager *Manager) Resolve(ctx context.Context, id ID, resource Resource) error {
	resource, set, err := manager.resources.GetOrSet(ctx, id, resource)
	if err != nil {
		return err
	}
	if !set {
		return ErrAlreadyResolved
	}
	return manager.resolutions.Write(ctx, id)
}

func (manager *Manager) Load(ctx context.Context, id ID) (Resolution, error) {
	resource, err := manager.resources.Get(ctx, id)
	if err != nil {
		if err == data.NotFoundErr {
			return Resolution{}, ErrNotResolved
		}
		return Resolution{}, fmt.Errorf("failed to load a resource: %s", err)
	}

	future, err := manager.futures.Get(ctx, id)
	if err != nil {
		if err != data.NotFoundErr {
			err = fmt.Errorf("failed to load a future: %s", err)
		}
		return Resolution{}, err
	}

	return Resolution{Future: future, Resource: resource}, nil
}

func (manager *Manager) Resolutions() pool.Reader[ID] {
	return manager.resolutions
}

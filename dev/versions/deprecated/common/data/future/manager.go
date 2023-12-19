package future

import (
	"context"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"kantoku/common/data"
	"kantoku/common/data/kv"
	"kantoku/common/data/pool"
)

type Manager struct {
	resolutions kv.Database[ID, Resolution]
	pool        pool.Pool[ID]
}

func NewManager(resolutions kv.Database[ID, Resolution], pool pool.Pool[ID]) *Manager {
	return &Manager{
		resolutions: resolutions,
		pool:        pool,
	}
}

func (manager *Manager) Make(_ context.Context) (ID, error) {
	id := uuid.New().String()
	return id, nil
}

func (manager *Manager) OK(ctx context.Context, id ID, resource Resource) error {
	return manager.Resolve(ctx, id, resource, OK)
}

func (manager *Manager) Fail(ctx context.Context, id ID, resource Resource) error {
	return manager.Resolve(ctx, id, resource, FAILURE)
}

func (manager *Manager) Resolve(ctx context.Context, id ID, resource Resource, status Status) error {
	res := Resolution{
		Future:   id,
		Resource: resource,
		Status:   status,
	}
	_, set, err := manager.resolutions.GetOrSet(ctx, id, res)
	if err != nil {
		return err
	}
	if !set {
		return ErrAlreadyResolved
	}

	err = manager.pool.Write(ctx, id)
	if err != nil {
		deletionErr := manager.resolutions.Del(ctx, id)

		message := fmt.Sprintf("failed to put resolution into the pool ('%s')", err)

		if deletionErr != nil {
			message += ", "
			message += fmt.Sprintf("failed to delete the resolution from a pool ('%s')", deletionErr)
		}

		return errors.New(message)
	}

	return nil
}

func (manager *Manager) Load(ctx context.Context, id ID) (Resolution, error) {
	res, err := manager.resolutions.Get(ctx, id)
	if err != nil {
		if err == data.NotFoundErr {
			return Resolution{}, ErrNotResolved
		}
		return Resolution{}, fmt.Errorf("failed to load a resource: %s", err)
	}

	return res, nil
}

func (manager *Manager) Resolutions() pool.Reader[ID] {
	return manager.pool
}

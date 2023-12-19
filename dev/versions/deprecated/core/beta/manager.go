package beta

import (
	"context"
	"fmt"
	"kantoku/common/data/bimap"
	"kantoku/common/data/deps"
	"kantoku/core/alpha"
)

type Manager struct {
	alphas       *alpha.Manager
	dependencies deps.Deps
	alpha2group  bimap.Bimap[string, string]
}

func (manager *Manager) Get(id string) Beta {
	return Beta{manager, id}
}

func (manager *Manager) Spawn(ctx context.Context, spec Spec) error {
	alpha, err := manager.alphas.New(ctx, spec.Data)
	if err != nil {
		return fmt.Errorf("failed to create a new alpha: %s", err)
	}

	id, err := manager.dependencies.NewGroup(ctx)
	if err != nil {
		return fmt.Errorf("failed to create a new group: %s", err)
	}

	if err := manager.alpha2group.Set(ctx, alpha.ID(), id); err != nil {

	}
}

package delay

import (
	"context"
	"kantoku/common/data/cron"
	"kantoku/unused/backend/framework/depot/deps"
	"time"
)

type Manager struct {
	cron cron.Cron
	deps deps.Deps
}

func NewManager(cron cron.Cron, deps deps.Deps) *Manager {
	return &Manager{
		cron: cron,
		deps: deps,
	}
}

func (manager *Manager) MakeDependency(ctx context.Context, at time.Time) (string, error) {
	dep, err := manager.deps.Make(ctx)
	if err != nil {
		return "", err
	}

	if err := manager.cron.Schedule(ctx, at, dep.ID); err != nil {
		return "", err
	}

	return dep.ID, nil
}

func (manager *Manager) Cron() cron.Cron {
	return manager.cron
}

func (manager *Manager) Deps() deps.Deps {
	return manager.deps
}

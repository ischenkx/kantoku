package depot

import (
	"context"
	"kantoku/common/db/kv"
	"kantoku/common/deps"
	"kantoku/core/l2"
	"log"
)

type Manager struct {
	deps          deps.DB
	group2task    kv.Database[string]
	taskScheduler l2.Scheduler
}

func NewManager(deps deps.DB, group2task kv.Database[string], taskScheduler l2.Scheduler) *Manager {
	return &Manager{
		deps:          deps,
		group2task:    group2task,
		taskScheduler: taskScheduler,
	}
}

func (manager *Manager) Schedule(ctx context.Context, task Task) error {
	group, err := manager.Deps().Make(ctx, task.Dependencies...)
	if err != nil {
		return err
	}

	if _, err := manager.group2task.Set(ctx, group.ID, task.ID); err != nil {
		return err
	}

	return nil
}

func (manager *Manager) Deps() deps.DB {
	return manager.deps
}

func (manager *Manager) Run(ctx context.Context) error {
	ready, err := manager.Deps().Ready(ctx)
	if err != nil {
		return err
	}

loop:
	for {
		select {
		case <-ctx.Done():
			break loop
		case id := <-ready:
			if err := manager.taskScheduler.Write(ctx, id); err != nil {
				log.Println("failed to schedule a task:", err)
				continue
			}
		}
	}

	return nil
}

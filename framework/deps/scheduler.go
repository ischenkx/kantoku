package deps

import (
	"context"
	"kantoku/common/deps"
	"kantoku/core/l2"
	"log"
)

type Scheduler struct {
	deps          deps.DB
	taskScheduler l2.Scheduler
}

func (s *Scheduler) Schedule(ctx context.Context, task Task) error {
	mapDeps := map[string]bool{}
	for _, dep := range task.Dependencies {
		mapDeps[dep] = false
	}

	return s.Deps().MakeGroup(ctx, deps.Group{
		ID:           task.ID,
		Dependencies: mapDeps,
	})
}

func (s *Scheduler) Deps() deps.DB {
	return s.deps
}

func (s *Scheduler) Run(ctx context.Context) {
	ready := s.Deps().Ready(ctx)

loop:
	for {
		select {
		case <-ctx.Done():
			break loop
		case wg := <-ready:
			if err := s.taskScheduler.Schedule(ctx, wg.ID); err != nil {
				log.Println("failed to schedule a task:", err)
				continue
			}
		}
	}
}

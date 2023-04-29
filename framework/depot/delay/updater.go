package delay

import (
	"context"
	"log"
)

type Updater struct {
	manager *Manager
}

func NewUpdater(manager *Manager) *Updater {
	return &Updater{manager: manager}
}

func (updater *Updater) Run(ctx context.Context) error {
	events, err := updater.manager.Cron().Events(ctx)
	if err != nil {
		return err
	}

updater:
	for {
		select {
		case <-ctx.Done():
			break updater
		case dependencyID := <-events:
			if err := updater.manager.Deps().Resolve(ctx, dependencyID); err != nil {
				log.Println("failed to resolve a dependency:", err)
			}
		}
	}

	return nil
}

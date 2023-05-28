package taskdep

import (
	"context"
	"kantoku/platform"
	"kantoku/unused/backend/framework"
	"log"
)

type Updater struct {
	events  platform.Broker
	manager *Manager
}

func NewUpdater(events platform.Broker, manager *Manager) *Updater {
	return &Updater{
		events:  events,
		manager: manager,
	}
}

func (updater *Updater) Run(ctx context.Context) error {
	events, err := updater.events.Listen(ctx, framework.EventTopic)
	if err != nil {
		return err
	}

updater:
	for {
		select {
		case <-ctx.Done():
			break updater
		case ev := <-events:
			if ev.Name == framework.SentOutputsEvent {
				id := string(ev.Data)
				if err := updater.manager.ResolveTask(ctx, id); err != nil {
					log.Println("failed to resolve a task:", err)
				}
			}
		}
	}

	return nil
}

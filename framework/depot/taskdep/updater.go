package taskdep

import (
	"context"
	"kantoku/core/event"
	"kantoku/core/task"
	"log"
)

type Updater struct {
	events  event.Bus
	manager *Manager
}

func NewUpdater(events event.Bus, manager *Manager) *Updater {
	return &Updater{
		events:  events,
		manager: manager,
	}
}

func (updater *Updater) Run(ctx context.Context) error {
	events, err := updater.events.Listen(ctx, task.EventTopic)
	if err != nil {
		return err
	}

updater:
	for {
		select {
		case <-ctx.Done():
			break updater
		case ev := <-events:
			if ev.Name == task.SentOutputsEvent {
				id := string(ev.Data)
				if err := updater.manager.ResolveTask(ctx, id); err != nil {
					log.Println("failed to resolve a task:", err)
				}
			}
		}
	}

	return nil
}

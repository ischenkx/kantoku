package status

import (
	"context"
	"kantoku/core/event"
	"kantoku/core/task"
	"log"
)

type Updater struct {
	bus event.Bus
	db  DB
}

func NewUpdater(bus event.Bus, db DB) *Updater {
	return &Updater{
		bus: bus,
		db:  db,
	}
}

func (updater *Updater) Run(ctx context.Context) {
	events, err := updater.bus.Listen(ctx, task.EventTopic, task.EventTopic)
	if err != nil {
		log.Println("failed to subscribe to task and task events:", err)
		return
	}

loop:
	for {
		select {
		case <-ctx.Done():
			break loop

		case ev := <-events:
			id := string(ev.Data)
			switch ev.Name {
			case task.ScheduledEvent:
				updater.update(ctx, id, Pending)
			case task.ReceivedEvent:
				updater.update(ctx, id, Executing)
			case task.ExecutedEvent:
				updater.update(ctx, id, Executed)
			case task.SentOutputsEvent:
				updater.update(ctx, id, Complete)
			}
		}
	}
}

func (updater *Updater) update(ctx context.Context, id string, status Status) {
	if err := updater.db.UpdateStatus(ctx, id, status); err != nil {
		log.Printf("failed to update the status: id = '%s', status = '%s'\n", id, status)
	}
}

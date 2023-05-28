package status

import (
	"context"
	"kantoku/common/data/kv"
	"kantoku/platform"
	"kantoku/unused/backend/framework"
	"log"
)

type Updater struct {
	bus platform.Broker
	db  kv.Database[string, Status]
}

func NewUpdater(bus platform.Broker, db kv.Database[string, Status]) *Updater {
	return &Updater{
		bus: bus,
		db:  db,
	}
}

func (updater *Updater) Run(ctx context.Context) {
	events, err := updater.bus.Listen(ctx, framework.EventTopic)
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
			case framework.ScheduledEvent:
				updater.update(ctx, id, Pending)
			case framework.ReceivedEvent:
				updater.update(ctx, id, Executing)
			case framework.ExecutedEvent:
				updater.update(ctx, id, Executed)
			case framework.SentOutputsEvent:
				updater.update(ctx, id, Complete)
			}
		}
	}
}

func (updater *Updater) update(ctx context.Context, id string, status Status) {
	if err := updater.db.Set(ctx, id, status); err != nil {
		log.Printf("failed to update the status: id = '%s', status = '%s'\n", id, status)
	}
}

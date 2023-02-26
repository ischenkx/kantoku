package status

import (
	"context"
	"kantoku/core/l0/event"
	"kantoku/core/l1"
	"kantoku/core/l2"
	"log"
)

type Updater struct {
	bus event.Bus
	db  DB
}

func (updater *Updater) Run(ctx context.Context) {
	events, err := updater.bus.Listen(ctx, l1.EventTopic, l2.EventTopic)
	if err != nil {
		log.Println("failed to subscribe to l1 and l2 events:", err)
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
			case l2.SentTaskEvent:
				updater.update(ctx, id, Pending)
			case l1.ReceivedTaskEvent:
				updater.update(ctx, id, Executing)
			case l1.ExecutedTaskEvent:
				updater.update(ctx, id, Executed)
			case l1.SentOutputsEvent:
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

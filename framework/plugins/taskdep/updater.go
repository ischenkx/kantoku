package taskdep

import (
	"context"
	"kantoku/kernel/platform"
	"log"
)

type Updater struct {
	events                platform.Broker
	sentOutputsEventTopic string
	manager               *Manager
}

func NewUpdater(events platform.Broker, manager *Manager, sentOutputsEventTopic string) *Updater {
	return &Updater{
		events:                events,
		manager:               manager,
		sentOutputsEventTopic: sentOutputsEventTopic,
	}
}

func (updater *Updater) Run(ctx context.Context) error {
	listener := updater.events.Listen()
	defer listener.Close(ctx)

	err := listener.Subscribe(ctx, updater.sentOutputsEventTopic)
	if err != nil {
		return err
	}

	channel, err := listener.Incoming(ctx)
	if err != nil {
		return err
	}
updater:
	for {
		select {
		case <-ctx.Done():
			break updater
		case ev := <-channel:
			id := string(ev.Data)
			if err := updater.manager.ResolveTask(ctx, id); err != nil {
				log.Println("failed to resolve a task:", err)
			}
		}
	}

	return nil
}

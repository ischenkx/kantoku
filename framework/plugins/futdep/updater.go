package futdep

import (
	"context"
	"kantoku/common/data/pool"
	"kantoku/framework/future"
	"log"
)

type Updater struct {
	resolvedFutures       pool.Reader[future.ID]
	sentOutputsEventTopic string
	manager               *Manager
}

func NewUpdater(manager *Manager, resolvedFutures pool.Reader[future.ID]) *Updater {
	return &Updater{
		resolvedFutures: resolvedFutures,
		manager:         manager,
	}
}

func (updater *Updater) Run(ctx context.Context) error {
	channel, err := updater.resolvedFutures.Read(ctx)
	if err != nil {
		return err
	}
updater:
	for {
		select {
		case <-ctx.Done():
			break updater
		case id := <-channel:
			if err := updater.manager.ResolveFuture(ctx, id); err != nil {
				log.Println("failed to resolve a future:", err)
			}
		}
	}

	return nil
}

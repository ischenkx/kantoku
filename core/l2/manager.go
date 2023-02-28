package l2

import (
	"context"
	"kantoku/common/db/kv"
	"kantoku/common/pool"
	"kantoku/core/l0/event"
	"kantoku/core/l1"
	"log"
)

type Manager[T Task] struct {
	scheduler Scheduler
	db        kv.Database[T]
	events    event.Bus
	pool      pool.Writer[l1.Task]
}

func NewManager[T Task](scheduler Scheduler, db kv.Database[T], events event.Bus, pool pool.Writer[l1.Task]) *Manager[T] {
	return &Manager[T]{
		scheduler: scheduler,
		db:        db,
		events:    events,
		pool:      pool,
	}
}

func (manager *Manager[T]) Scheduler() Scheduler {
	return manager.scheduler
}

func (manager *Manager[T]) DB() kv.Database[T] {
	return manager.db
}

func (manager *Manager[T]) Run(ctx context.Context) {
	if err := manager.processPendingTasks(ctx); err != nil {
		log.Println("failed to process pending tasks:", err)
	}
}

func (manager *Manager[T]) processPendingTasks(ctx context.Context) error {
	pending, err := manager.Scheduler().Read(ctx)
	if err != nil {
		return err
	}
loop:
	for {
		select {
		case <-ctx.Done():
			break loop
		case id := <-pending:
			task, err := manager.db.Get(ctx, id)
			if err != nil {
				log.Println("failed to get a task from the database:", err)
				continue
			}

			l1task, err := task.AsL1(ctx)
			if err != nil {
				log.Println("failed to convert a task to l1:", err)
				continue
			}

			if err := manager.pool.Write(ctx, l1task); err != nil {
				log.Println("failed to send task to the pool:", err)
				continue
			}
			manager.publish(ctx, SentTaskEvent, []byte(l1task.ID))
		}
	}

	return nil
}

func (manager *Manager[T]) publish(ctx context.Context, name string, data []byte) {
	err := manager.events.Publish(ctx, event.Event{
		Topic: EventTopic,
		Name:  name,
		Data:  data,
	})
	if err != nil {
		log.Println("failed to publish an event:", err)
	}
}

package l2

import (
	"context"
	"kantoku/common/db"
	event2 "kantoku/core/l0/event"
	l12 "kantoku/core/l1"
	"log"
)

type Manager[T Task] struct {
	scheduler Scheduler
	db        db.KV[T]
	events    event2.Bus
	pool      l12.PoolWriter[l12.Task]
}

func (manager *Manager[T]) Scheduler() Scheduler {
	return manager.scheduler
}

func (manager *Manager[T]) DB() db.KV[T] {
	return manager.db
}

func (manager *Manager[T]) Run(ctx context.Context) {
	go manager.processPendingTasks(ctx)
}

func (manager *Manager[T]) processPendingTasks(ctx context.Context) {
	pending := manager.Scheduler().Pending(ctx)

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

			if err := manager.pool.Put(ctx, l1task); err != nil {
				log.Println("failed to send task to the pool:", err)
				continue
			}
			manager.publish(ctx, SentTaskEvent, []byte(l1task.ID))
		}
	}
}

func (manager *Manager[T]) publish(ctx context.Context, name string, data []byte) {
	err := manager.events.Publish(ctx, event2.Event{
		Topic: EventTopic,
		Name:  name,
		Data:  data,
	})
	if err != nil {
		log.Println("failed to publish an event:", err)
	}
}

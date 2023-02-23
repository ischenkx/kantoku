package l2

import (
	"context"
	"kantoku/l0/event"
	"kantoku/l1"
	"log"
)

type Scheduler[T Task] struct {
	db     DB[T]
	events event.Bus
	pool   l1.PoolWriter[l1.Task]
}

func (manager *Scheduler[T]) DB() DB[T] {
	return manager.db
}

func (manager *Scheduler[T]) Run(ctx context.Context) {
	go manager.processPendingTasks(ctx)
}

func (manager *Scheduler[T]) processPendingTasks(ctx context.Context) {
	pending := manager.DB().Pending(ctx)

loop:
	for {
		select {
		case <-ctx.Done():
			break loop
		case task := <-pending:
			l1Task, err := task.ToL1(ctx)
			if err != nil {
				log.Println("failed to cast a pending task to an l1 task:", err)
				continue
			}

			if err := manager.pool.Put(ctx, l1Task); err != nil {
				log.Println("failed to send task to the pool:", err)
				continue
			}

			manager.publish(ctx, SentTaskEvent, []byte(l1Task.ID))
		}
	}
}

func (manager *Scheduler[T]) publish(ctx context.Context, name string, data []byte) {
	err := manager.events.Publish(ctx, event.Event{
		Topic: EventTopic,
		Name:  name,
		Data:  data,
	})
	if err != nil {
		log.Println("failed to publish an event:", err)
	}
}

package l2

import (
	"context"
	"kantoku/common/data/kv"
	"kantoku/common/pool"
	"kantoku/core/l0/event"
	"kantoku/core/l1"
	"log"
)

type L2[T Task] struct {
	db     kv.Reader[T]
	events event.Bus
	pool   pool.Writer[l1.Task]
}

func New[T Task](db kv.Reader[T], events event.Bus, pool pool.Writer[l1.Task]) *L2[T] {
	return &L2[T]{
		db:     db,
		events: events,
		pool:   pool,
	}
}

func (l2 *L2[T]) Run(ctx context.Context, id string) error {
	task, err := l2.db.Get(ctx, id)
	if err != nil {
		return err
	}

	l1task, err := task.L1(ctx)
	if err != nil {
		return err
	}

	if err := l2.pool.Write(ctx, l1task); err != nil {
		return err
	}
	l2.publish(ctx, SentTaskEvent, []byte(l1task.ID))

	return nil
}

func (l2 *L2[T]) publish(ctx context.Context, name string, data []byte) {
	err := l2.events.Publish(ctx, event.Event{
		Topic: EventTopic,
		Name:  name,
		Data:  data,
	})
	if err != nil {
		log.Println("failed to publish an event:", err)
	}
}

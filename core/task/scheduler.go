package task

import (
	"context"
	"kantoku/common/data/pool"
	"kantoku/core/event"
	"log"
)

type SchedulerInputs[InputType AbstractTask] pool.Writer[InputType]

type Scheduler[InputType AbstractTask] struct {
	inputs SchedulerInputs[InputType]
	events event.Bus
}

func NewScheduler[InputType AbstractTask](
	inputs SchedulerInputs[InputType],
	events event.Bus) *Scheduler[InputType] {
	return &Scheduler[InputType]{
		inputs: inputs,
		events: events,
	}
}

func (scheduler *Scheduler[InputType]) Schedule(ctx context.Context, input InputType) error {
	if err := scheduler.inputs.Write(ctx, input); err != nil {
		return err
	}
	scheduler.publish(ctx, ScheduledEvent, []byte(input.ID()))

	return nil
}

func (scheduler *Scheduler[InputType]) publish(ctx context.Context, name string, data []byte) {
	err := scheduler.events.Publish(ctx, event.Event{
		Topic: EventTopic,
		Name:  name,
		Data:  data,
	})
	if err != nil {
		log.Println("failed to publish an event:", err)
	}
}

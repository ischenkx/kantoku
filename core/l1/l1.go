package l1

import (
	"context"
	"kantoku/common/pool"
	event2 "kantoku/core/l0/event"
	"log"
)

type L1 struct {
	inputs   pool.Reader[Task]
	outputs  pool.Writer[Result]
	executor Executor
	events   event2.Publisher
}

func New(inputs pool.Reader[Task], outputs pool.Writer[Result], executor Executor, events event2.Publisher) *L1 {
	return &L1{
		inputs:   inputs,
		outputs:  outputs,
		executor: executor,
		events:   events,
	}
}

func (l1 *L1) Run(ctx context.Context) error {
	// todo: inputs are not closed explicitly in this function
	// todo: maybe I should create another interface for a
	// todo: closeable read-only channel (though sounds kinda broken by design)
	inputs, err := l1.inputs.Read(ctx)
	if err != nil {
		return err
	}

loop:
	for {
		select {
		case <-ctx.Done():
			break loop
		case task := <-inputs:
			// todo: remove implicit topic assignment
			// todo: I should leave the logic of determining of what topic
			// todo: to use to the Publisher provided by a user
			l1.sendEvent(ctx, ReceivedTaskEvent, EventTopic, []byte(task.ID))

			result, err := l1.executor.Execute(task)
			if err != nil {
				log.Printf("failed to execute a task (id = '%s'): %s\n", task.ID, err)
				result = Result{
					TaskID: task.ID,
					Data:   []byte(err.Error()),
					Status: FAILURE,
				}
			}

			l1.sendEvent(ctx, ExecutedTaskEvent, EventTopic, []byte(task.ID))

			err = l1.outputs.Write(ctx, result)
			if err != nil {
				log.Printf("failed to save the output of a task (id = '%s'): %s\n", task.ID, err)
			}

			l1.sendEvent(ctx, SentOutputsEvent, EventTopic, []byte(task.ID))
		}
	}
	return nil
}

func (l1 *L1) sendEvent(ctx context.Context, name string, topic string, data []byte) {
	err := l1.events.Publish(ctx, event2.Event{
		Name:  name,
		Topic: topic,
		Data:  data,
	})
	if err != nil {
		log.Printf("failed to send the event: %s", err)
	}
}

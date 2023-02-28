package l1

import (
	"context"
	"kantoku/common/pool"
	event2 "kantoku/core/l0/event"
	"log"
)

type Runner struct {
	inputs   pool.Reader[Task]
	outputs  pool.Writer[Result]
	executor Executor
	events   event2.Publisher
}

func NewRunner(inputs pool.Reader[Task], outputs pool.Writer[Result], executor Executor, events event2.Publisher) *Runner {
	return &Runner{
		inputs:   inputs,
		outputs:  outputs,
		executor: executor,
		events:   events,
	}
}

func (runner *Runner) Run(ctx context.Context) error {
	// todo: inputs are not closed explicitly in this function
	// todo: maybe I should create another interface for a
	// todo: closeable read-only channel (though sounds kinda broken by design)
	inputs, err := runner.inputs.Read(ctx)
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
			runner.sendEvent(ctx, ReceivedTaskEvent, EventTopic, []byte(task.ID))

			result, err := runner.executor.Execute(task)
			if err != nil {
				log.Printf("failed to execute a task (id = '%s'): %s\n", task.ID, err)
				result = Result{
					TaskID: task.ID,
					Data:   []byte(err.Error()),
					Status: FAILURE,
				}
			}

			runner.sendEvent(ctx, ExecutedTaskEvent, EventTopic, []byte(task.ID))

			err = runner.outputs.Write(ctx, result)
			if err != nil {
				log.Printf("failed to save the output of a task (id = '%s'): %s\n", task.ID, err)
			}

			runner.sendEvent(ctx, SentOutputsEvent, EventTopic, []byte(task.ID))
		}
	}
	return nil
}

func (runner *Runner) sendEvent(ctx context.Context, name string, topic string, data []byte) {
	err := runner.events.Publish(ctx, event2.Event{
		Name:  name,
		Topic: topic,
		Data:  data,
	})
	if err != nil {
		log.Printf("failed to send the event: %s", err)
	}
}

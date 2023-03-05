package task

import (
	"context"
	"kantoku/common/pool"
	"kantoku/core/event"
	"log"
)

type PipelineInputs[T AbstractTask] pool.Reader[T]
type PipelineOutputs pool.Writer[Result]

type Pipeline[InputType AbstractTask] struct {
	inputs   PipelineInputs[InputType]
	outputs  PipelineOutputs
	executor Executor[InputType]
	events   event.Publisher
}

func NewPipeline[InputType AbstractTask](
	inputs PipelineInputs[InputType],
	outputs PipelineOutputs,
	executor Executor[InputType],
	events event.Publisher) *Pipeline[InputType] {
	return &Pipeline[InputType]{
		inputs:   inputs,
		outputs:  outputs,
		executor: executor,
		events:   events,
	}
}

func (pipeline *Pipeline[InputType]) Run(ctx context.Context) error {
	// todo: inputs are not closed explicitly in this function
	// todo: maybe I should create another interface for a
	// todo: closeable read-only channel (though sounds kinda broken by design)
	inputs, err := pipeline.inputs.Read(ctx)
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
			id := task.ID(ctx)
			pipeline.sendEvent(ctx, ReceivedEvent, EventTopic, []byte(id))

			result, err := pipeline.executor.Execute(ctx, task)
			if err != nil {
				log.Printf("failed to execute a task (id = '%s'): %s\n", id, err)
				result = Result{
					TaskID: id,
					Data:   []byte(err.Error()),
					Status: FAILURE,
				}
			}

			pipeline.sendEvent(ctx, ExecutedEvent, EventTopic, []byte(id))

			err = pipeline.outputs.Write(ctx, result)
			if err != nil {
				log.Printf("failed to save the output of a task (id = '%s'): %s\n", id, err)
			}

			pipeline.sendEvent(ctx, SentOutputsEvent, EventTopic, []byte(id))
		}
	}
	return nil
}

func (pipeline *Pipeline[InputType]) sendEvent(ctx context.Context, name string, topic string, data []byte) {
	err := pipeline.events.Publish(ctx, event.Event{
		Name:  name,
		Topic: topic,
		Data:  data,
	})
	if err != nil {
		log.Printf("failed to send the event: %s", err)
	}
}

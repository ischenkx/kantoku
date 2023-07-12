package evexec

import (
	"context"
	"kantoku/backend/executor"
	"kantoku/common/data/pool"
	"kantoku/kernel/platform"
	"log"
)

// Executor is an implementation of the Task Events Protocol
type Executor[Task platform.Task] struct {
	runner        executor.Runner[Task, []byte]
	platform      platform.Platform[Task]
	topicResolver TopicResolver
}

// Run - starts the queue processing pipeline (synchronous)
//
// In order to process multiple tasks concurrently simply call this function in multiple Goroutines.
func (e *Executor[Task]) Run(ctx context.Context) error {
	return pool.ReadAutoCommit[string](ctx, e.platform.Inputs(), func(ctx context.Context, id string) error {
		// TODO: split this code into several methods, so I can get rid of all those nasty else's
		e.emit(ctx, platform.Event{Name: ReceivedEvent, Data: []byte(id)})
		task, err := e.platform.DB().Get(ctx, id)
		if err != nil {
			message, err := ErrorMessage{TaskID: id, Message: err.Error()}.Encode()
			if err != nil {
				log.Println("failed to generate an error message:", err)
			} else {
				e.emit(ctx, platform.Event{Name: ErrorEvent, Data: message})
			}
			return nil
		}

		output, err := e.runner.Run(ctx, task)
		e.emit(ctx, platform.Event{Name: ExecutedEvent, Data: []byte(task.ID())})

		result := platform.Result{TaskID: task.ID()}
		if err != nil {
			result.Data = []byte(err.Error())
			result.Status = platform.FAILURE
		} else {
			result.Data = output
			result.Status = platform.OK
		}

		err = e.platform.Outputs().Set(ctx, result.TaskID, result)
		if err != nil {
			message, err := ErrorMessage{TaskID: result.TaskID, Message: err.Error()}.Encode()
			if err != nil {
				log.Println("failed to generate an error message:", err)
			} else {
				e.emit(ctx, platform.Event{Name: ErrorEvent, Data: message})
			}
		} else {
			e.emit(ctx, platform.Event{Name: SentOutputsEvent, Data: []byte(result.TaskID)})
		}

		return nil
	})
}

func (e *Executor[Task]) emit(ctx context.Context, event platform.Event) {
	var err error
	event.Topic, err = e.topicResolver.Resolve(event.Name)

	if err != nil {
		log.Printf("failed to resolve the topic name (event='%s'): %s\n", event.Name, err)
		return
	}

	if err := e.platform.Broker().Publish(ctx, event); err != nil {
		log.Println("failed to publish event:", err)
	}
}

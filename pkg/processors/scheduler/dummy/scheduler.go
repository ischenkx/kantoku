package dummy

import (
	"context"
	"fmt"
	"github.com/ischenkx/kantoku/pkg/system"
	event2 "github.com/ischenkx/kantoku/pkg/system/kernel/event"
	"log/slog"
)

const QueueName = "scheduler"

type Processor struct {
	system system.AbstractSystem
}

func NewProcessor(sys system.AbstractSystem) *Processor {
	return &Processor{system: sys}
}

func (processor *Processor) Process(ctx context.Context) error {
	channel, err := processor.system.Events().Consume(ctx, event2.Queue{
		Name:   QueueName,
		Topics: []string{system.TaskNewEvent},
	})
	if err != nil {
		return fmt.Errorf("failed to read events: %w", err)
	}

	for {
		select {
		case <-ctx.Done():
			return nil
		case ev := <-channel:
			taskId := string(ev.Data)

			slog.Info("received a new task",
				slog.String("id", taskId))

			tasks, err := processor.system.Tasks().Load(ctx, taskId)
			if err != nil {
				// TODO: use transactional events
				slog.Error("failed to load a task",
					slog.String("id", taskId),
					slog.String("error", err.Error()),
				)
				continue
			}

			// this must be a valid action
			task := tasks[0]

			err = processor.system.
				Events().
				Publish(ctx, event2.New(system.TaskReadyEvent, []byte(task.ID)))
			if err != nil {
				// TODO: use transactional events
				slog.Error("failed to publish an event",
					slog.String("id", taskId),
					slog.String("error", err.Error()),
				)
				continue
			}
		}
	}
}

package executor

import (
	"context"
	"fmt"
	codec "github.com/ischenkx/kantoku/pkg/common/data/codec"
	"github.com/ischenkx/kantoku/pkg/system"
	event2 "github.com/ischenkx/kantoku/pkg/system/kernel/event"
	"log/slog"
)

const QueueName = "executor"

type Processor struct {
	system      system.AbstractSystem
	resultCodec codec.Codec[Result, []byte]
	localQueue  string
	executor    Executor
}

func NewProcessor(
	system system.AbstractSystem,
	executor Executor,
	localQueue string,
	resultCodec codec.Codec[Result, []byte]) *Processor {
	return &Processor{
		system:      system,
		resultCodec: resultCodec,
		executor:    executor,
		localQueue:  localQueue,
	}
}

func (processor *Processor) Process(ctx context.Context) error {
	localContext, cancel := context.WithCancel(ctx)
	defer cancel()

	readyTaskEvents, err := processor.system.Events().Consume(localContext, event2.Queue{
		Name:   QueueName,
		Topics: []string{system.TaskReadyEvent},
	})
	if err != nil {
		return fmt.Errorf("failed to read events: %w", err)
	}

	controller := newExecutionController(
		processor.system,
		processor.executor,
		processor.localQueue,
		processor.resultCodec,
	)
	if err := controller.start(localContext); err != nil {
		return fmt.Errorf("failed to start the controller: %w", err)
	}

	for {
		select {
		case <-ctx.Done():
			return nil
		case ev := <-readyTaskEvents:
			taskId := string(ev.Data)

			if err := controller.processReadyTask(ctx, taskId); err != nil {
				slog.Error("failed to process a ready task",
					slog.String("id", taskId),
					slog.String("error", err.Error()))
			}
		}
	}
}

package status

import (
	"context"
	"fmt"
	codec "github.com/ischenkx/kantoku/pkg/common/data/codec"
	"github.com/ischenkx/kantoku/pkg/common/data/record"
	"github.com/ischenkx/kantoku/pkg/processors/executor"
	"github.com/ischenkx/kantoku/pkg/system"
	event2 "github.com/ischenkx/kantoku/pkg/system/kernel/event"
	"github.com/ischenkx/kantoku/pkg/system/kernel/task"
	"log/slog"
)

var QueueName = "scheduler"

type Processor struct {
	system      system.AbstractSystem
	resultCodec codec.Codec[executor.Result, []byte]
}

func NewProcessor(sys *system.System, resultCodec codec.Codec[executor.Result, []byte]) *Processor {
	return &Processor{
		system:      sys,
		resultCodec: resultCodec,
	}
}

func (processor *Processor) Process(ctx context.Context) error {
	events, err := processor.system.
		Events().
		Consume(ctx,
			event2.Queue{
				Name: QueueName,
				Topics: []string{
					system.TaskNewEvent,
					system.TaskReadyEvent,
					system.TaskReceivedEvent,
					system.TaskFinishedEvent,
					system.TaskCancelledEvent,
				},
			})
	if err != nil {
		return fmt.Errorf("failed to consumer events: %w", err)
	}

processor:
	for {
		select {
		case <-ctx.Done():
			break processor
		case ev := <-events:
			slog.Debug("received an event",
				slog.String("event", ev.Topic),
				slog.Any("data", string(ev.Data)))
			processor.processEvent(ctx, ev)
		}
	}

	return nil
}

func (processor *Processor) processEvent(ctx context.Context, ev event2.Event) {
	switch ev.Topic {
	case system.TaskNewEvent, system.TaskReadyEvent, system.TaskReceivedEvent, system.TaskCancelledEvent:
		taskId := string(ev.Data)
		newStatus := processor.event2status(ev.Topic)
		if err := processor.updateStatus(ctx, taskId, newStatus); err != nil {
			slog.Error("failed to update status",
				slog.String("id", taskId),
				slog.String("status", newStatus),
				slog.String("error", err.Error()))
			break
		}

	case system.TaskFinishedEvent:
		result, err := processor.resultCodec.Decode(ev.Data)
		if err != nil {
			slog.Error("failed to decode the result",
				slog.String("data", string(ev.Data)),
				slog.String("error", err.Error()))
			break
		}

		newStatus := task.OkStatus
		if result.Status == executor.Failed {
			newStatus = task.FailedStatus
		}

		if err := processor.updateStatus(ctx, result.TaskID, newStatus); err != nil {
			slog.Error("failed to update status",
				slog.String("id", result.TaskID),
				slog.String("status", newStatus),
				slog.String("error", err.Error()))
			break
		}

		if err := processor.saveResultData(ctx, result); err != nil {
			slog.Error("failed to save a result",
				slog.String("id", result.TaskID),
				slog.String("error", newStatus))
			break
		}

	default:
		slog.Error("unexpected event topic",
			slog.String("topic", ev.Topic))
	}
}

func (processor *Processor) event2status(topic string) string {
	switch topic {
	case system.TaskNewEvent:
		return task.InitializedStatus
	case system.TaskReadyEvent:
		return task.ReadyStatus
	case system.TaskReceivedEvent:
		return task.ReceivedStatus
	case system.TaskCancelledEvent:
		return task.CancelledStatus
	default:
		return ""
	}
}

func (processor *Processor) updateStatus(ctx context.Context, id string, status string) error {
	err := processor.system.
		Info().
		Filter(record.R{system.InfoTaskID: id}).
		Update(ctx, record.R{"status": status}, record.R{system.InfoTaskID: id})
	if err != nil {
		return fmt.Errorf("failed to update records: %w", err)
	}

	return nil
}

func (processor *Processor) saveResultData(ctx context.Context, result executor.Result) error {
	err := processor.system.
		Info().
		Filter(record.R{"id": result.TaskID}).
		Update(ctx, record.R{"result": result.Data}, record.R{"id": result.TaskID})
	if err != nil {
		return fmt.Errorf("failed to update records: %w", err)
	}

	return nil
}

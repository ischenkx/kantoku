package status

import (
	"context"
	"fmt"
	"github.com/ischenkx/kantoku/pkg/common/broker"
	codec "github.com/ischenkx/kantoku/pkg/common/data/codec"
	"github.com/ischenkx/kantoku/pkg/common/data/record"
	"github.com/ischenkx/kantoku/pkg/common/data/record/ops"
	"github.com/ischenkx/kantoku/pkg/common/service"
	"github.com/ischenkx/kantoku/pkg/core/event"
	"github.com/ischenkx/kantoku/pkg/core/services/executor"
	"github.com/ischenkx/kantoku/pkg/core/task"

	"github.com/ischenkx/kantoku/pkg/core/system"
	"github.com/ischenkx/kantoku/pkg/core/system/events"
	"log/slog"
	"time"
)

var QueueName = "scheduler"

type Service struct {
	System      system.AbstractSystem
	ResultCodec codec.Codec[executor.Result, []byte]

	service.Core
}

func (srvc *Service) Run(ctx context.Context) error {
	evs, err := srvc.System.
		Events().
		Consume(ctx,
			broker.TopicsInfo{
				Group: QueueName,
				Topics: []string{
					events.OnTask.Created,
					events.OnTask.Ready,
					events.OnTask.Received,
					events.OnTask.Finished,
					events.OnTask.Cancelled,
				},
			})
	if err != nil {
		return fmt.Errorf("failed to consumer events: %w", err)
	}

	broker.Processor[event.Event]{
		Handler: func(ctx context.Context, ev event.Event) error {
			if err := srvc.processEvent(ctx, ev); err != nil {
				return err
			}

			return nil
		},
		ErrorHandler: func(ctx context.Context, ev event.Event, err error) {
			srvc.Logger().Error("processing failed",
				slog.String("error", err.Error()))
		},
	}.Process(ctx, evs)

	return nil
}

func (srvc *Service) processEvent(ctx context.Context, ev event.Event) error {
	switch ev.Topic {
	case events.OnTask.Created,
		events.OnTask.Ready,
		events.OnTask.Received,
		events.OnTask.Cancelled:

		taskId := string(ev.Data)
		newStatus := srvc.event2status(ev.Topic)
		if err := srvc.updateStatus(ctx, taskId, newStatus, ""); err != nil {
			return fmt.Errorf("failed to update status (task_id='%s' status='%s'): %w",
				taskId,
				newStatus,
				err)
		}
	case events.OnTask.Finished:
		result, err := srvc.ResultCodec.Decode(ev.Data)
		if err != nil {
			return fmt.Errorf("failed to decode the result: %w", err)
		}

		newStatus := task.FinishedStatus

		if err := srvc.updateStatus(ctx, result.TaskID, newStatus, string(result.Status)); err != nil {
			return fmt.Errorf("failed to update status (task_id='%s' status='%s'): %w",
				result.TaskID,
				newStatus,
				err)
		}

		if err := srvc.saveResultData(ctx, result); err != nil {
			return fmt.Errorf("failed to save result (task_id='%s'): %w",
				result.TaskID,
				err)
		}
	default:
		srvc.Logger().Error("unexpected event topic",
			slog.String("topic", ev.Topic))
	}

	return nil
}

func (srvc *Service) event2status(topic string) string {
	switch topic {
	case events.OnTask.Created:
		return task.InitializedStatus
	case events.OnTask.Ready:
		return task.ReadyStatus
	case events.OnTask.Received:
		return task.ReceivedStatus
	case events.OnTask.Cancelled:
		return task.CancelledStatus
	default:
		return ""
	}
}

func (srvc *Service) status2precedingStatuses(status string) []any {
	switch status {
	case task.InitializedStatus:
		return []any{nil}
	case task.ReadyStatus:
		return []any{task.InitializedStatus}
	case task.ReceivedStatus:
		return []any{task.InitializedStatus, task.ReadyStatus}
	case task.FinishedStatus:
		return []any{task.InitializedStatus, task.ReadyStatus, task.ReceivedStatus}
	case task.CancelledStatus:
		return []any{}
	default:
		return []any{}
	}
}

func (srvc *Service) updateStatus(ctx context.Context, id string, status, subStatus string) error {
	now := time.Now().Unix()

	err := srvc.System.
		Tasks().
		Filter(record.R{"id": id, "info.status": ops.In[any](srvc.status2precedingStatuses(status)...)}).
		Update(ctx,
			record.R{"info": record.R{"status": status, "sub_status": subStatus, "updated_at": now}},
			nil,
		)
	if err != nil {
		return fmt.Errorf("failed to update records: %w", err)
	}

	return nil
}

func (srvc *Service) saveResultData(ctx context.Context, result executor.Result) error {
	err := srvc.System.
		Tasks().
		Filter(record.R{"id": result.TaskID}).
		Update(ctx, record.R{"info": record.R{"result": result.Data}}, nil)
	if err != nil {
		return fmt.Errorf("failed to update records: %w", err)
	}

	return nil
}

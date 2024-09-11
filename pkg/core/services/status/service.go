package status

import (
	"context"
	"fmt"
	codec "github.com/ischenkx/kantoku/pkg/common/data/codec"
	"github.com/ischenkx/kantoku/pkg/common/service"
	"github.com/ischenkx/kantoku/pkg/common/transport/broker"
	"github.com/ischenkx/kantoku/pkg/core"
	"github.com/ischenkx/kantoku/pkg/core/services/executor"
	"log/slog"
	"time"
)

var QueueName = "status"

type Service struct {
	System      core.AbstractSystem
	ResultCodec codec.Codec[executor.Result, []byte]

	service.Core
}

func (srvc *Service) Run(ctx context.Context) error {
	evs, err := srvc.System.
		Events().
		Consume(ctx,
			[]string{
				core.OnTask.Created,
				core.OnTask.Ready,
				core.OnTask.Received,
				core.OnTask.Finished,
				core.OnTask.Cancelled,
			},
			broker.ConsumerSettings{
				Group:                QueueName,
				InitializationPolicy: broker.OldestOffset,
			},
		)
	if err != nil {
		return fmt.Errorf("failed to consumer events: %w", err)
	}

	broker.Processor[core.Event]{
		Handler: func(ctx context.Context, ev core.Event) error {
			if err := srvc.processEvent(ctx, ev); err != nil {
				return err
			}

			return nil
		},
		ErrorHandler: func(ctx context.Context, ev core.Event, err error) {
			srvc.Logger().Error("processing failed",
				slog.String("error", err.Error()))
		},
	}.Process(ctx, evs)

	return nil
}

func (srvc *Service) processEvent(ctx context.Context, ev core.Event) error {
	switch ev.Topic {
	case core.OnTask.Created,
		core.OnTask.Ready,
		core.OnTask.Received,
		core.OnTask.Cancelled:

		taskId := string(ev.Data)
		newStatus := srvc.event2status(ev.Topic)
		if err := srvc.updateStatus(ctx, taskId, newStatus, ""); err != nil {
			return fmt.Errorf("failed to update status (task_id='%s' status='%s'): %w",
				taskId,
				newStatus,
				err)
		}
	case core.OnTask.Finished:
		result, err := srvc.ResultCodec.Decode(ev.Data)
		if err != nil {
			return fmt.Errorf("failed to decode the result: %w", err)
		}

		newStatus := core.TaskStatuses.Finished

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
	case core.OnTask.Created:
		return core.TaskStatuses.Initialized
	case core.OnTask.Ready:
		return core.TaskStatuses.Ready
	case core.OnTask.Received:
		return core.TaskStatuses.Received
	case core.OnTask.Cancelled:
		return core.TaskStatuses.Cancelled
	default:
		return ""
	}
}

func (srvc *Service) status2precedingStatuses(status string) []any {
	switch status {
	case core.TaskStatuses.Initialized:
		return []any{nil}
	case core.TaskStatuses.Ready:
		return []any{core.TaskStatuses.Initialized}
	case core.TaskStatuses.Received:
		return []any{core.TaskStatuses.Initialized, core.TaskStatuses.Ready}
	case core.TaskStatuses.Finished:
		return []any{core.TaskStatuses.Initialized, core.TaskStatuses.Ready, core.TaskStatuses.Received}
	case core.TaskStatuses.Cancelled:
		return []any{}
	default:
		return []any{}
	}
}

func (srvc *Service) updateStatus(ctx context.Context, id string, status, subStatus string) error {
	now := time.Now().Unix()

	_, err := srvc.System.
		Tasks().
		UpdateWithProperties(
			ctx,
			map[string][]any{
				"id":          {id},
				"info.status": srvc.status2precedingStatuses(status),
			},
			map[string]any{
				"info.status":     status,
				"info.sub_status": subStatus,
				"info.updated_at": now,
			},
		)
	if err != nil {
		return fmt.Errorf("failed to update records: %w", err)
	}

	return nil
}

func (srvc *Service) saveResultData(ctx context.Context, result executor.Result) error {
	_, err := srvc.System.
		Tasks().
		UpdateWithProperties(
			ctx,
			map[string][]any{
				"id": {result.TaskID},
			},
			map[string]any{
				"info.result": string(result.Data),
			},
		)
	if err != nil {
		return fmt.Errorf("failed to update records: %w", err)
	}

	return nil
}

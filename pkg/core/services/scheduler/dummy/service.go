package dummy

import (
	"context"
	"fmt"
	"github.com/ischenkx/kantoku/pkg/common/service"
	"github.com/ischenkx/kantoku/pkg/common/transport/broker"
	"github.com/ischenkx/kantoku/pkg/core"
	"log/slog"
)

const QueueName = "scheduler"

type Service struct {
	System core.AbstractSystem
	service.Core
}

func (srvc *Service) Run(ctx context.Context) error {
	channel, err := srvc.System.Events().Consume(ctx, []string{core.OnTask.Created}, broker.ConsumerSettings{
		Group: QueueName,
	})
	if err != nil {
		return fmt.Errorf("failed to read events: %w", err)
	}

	broker.Processor[core.Event]{
		Handler: func(ctx context.Context, ev core.Event) error {
			taskId := string(ev.Data)

			srvc.Logger().Debug("new task",
				slog.String("id", taskId))

			t, err := srvc.System.Task(ctx, taskId)
			if err != nil {
				return fmt.Errorf("failed to load task (id='%s'): %w", taskId, err)
			}

			err = srvc.System.Events().Send(ctx, core.NewEvent(core.OnTask.Ready, []byte(t.ID)))
			if err != nil {
				return fmt.Errorf("failed to publish an event (taskId='%s'): %w", taskId, err)
			}

			return nil
		},
		ErrorHandler: func(ctx context.Context, ev core.Event, err error) {
			taskId := string(ev.Data)

			srvc.Logger().Error("failed to schedule task",
				slog.String("task_id", taskId),
				slog.String("error", err.Error()),
			)
		},
	}.Process(ctx, channel)

	return nil
}

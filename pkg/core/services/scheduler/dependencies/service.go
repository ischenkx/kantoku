package dependencies

import (
	"context"
	"fmt"
	"github.com/ischenkx/kantoku/pkg/common/service"
	"github.com/ischenkx/kantoku/pkg/common/transport/broker"
	"github.com/ischenkx/kantoku/pkg/core"
	"github.com/ischenkx/kantoku/pkg/core/services/scheduler/dependencies/manager"
	"golang.org/x/sync/errgroup"
	"log/slog"
)

var QueueName = "dependencies.simple"

type Service struct {
	System  core.AbstractSystem
	Manager *manager.Manager

	service.Core
}

func (srvc *Service) Run(ctx context.Context) error {
	g, ctx := errgroup.WithContext(ctx)
	g.Go(func() error {
		srvc.Logger().Info("processing incoming tasks...")
		return srvc.processNewTasks(ctx)
	})

	g.Go(func() error {
		srvc.Logger().Info("processing ready tasks...")
		return srvc.processReadyTasks(ctx)
	})

	if err := g.Wait(); err != nil {
		return err
	}

	return nil
}

func (srvc *Service) processNewTasks(ctx context.Context) error {
	channel, err := srvc.System.Events().Consume(
		ctx,
		[]string{core.OnTask.Created},
		broker.ConsumerSettings{Group: QueueName},
	)
	if err != nil {
		return fmt.Errorf("failed to read events: %w", err)
	}

	broker.Processor[core.Event]{
		Handler: func(ctx context.Context, ev core.Event) error {
			taskId := string(ev.Data)
			srvc.Logger().Debug("new task",
				slog.String("id", taskId))

			if err := srvc.Manager.Register(ctx, taskId); err != nil {
				srvc.Logger().Error("failed to process a created task",
					slog.String("task_id", taskId),
					slog.String("error", err.Error()))
				//return fmt.Errorf("failed to register a task (id='%s'): %w", taskId, err)
			}

			return nil
		},
		ErrorHandler: func(ctx context.Context, ev core.Event, err error) {
			//taskId := string(ev.Data)

		},
	}.Process(ctx, channel)

	return nil
}

func (srvc *Service) processReadyTasks(ctx context.Context) error {
	channel, err := srvc.Manager.Ready(ctx)
	if err != nil {
		return err
	}

	for {
		select {
		case <-ctx.Done():
			return nil
		case taskId := <-channel:
			srvc.Logger().Debug("ready task",
				slog.String("id", taskId))
			err := srvc.System.Events().Send(ctx, core.NewEvent(core.OnTask.Ready, []byte(taskId)))
			if err != nil {
				srvc.Logger().Error("failed to publish an event",
					slog.String("id", taskId),
					slog.String("event", core.OnTask.Ready),
					slog.String("error", err.Error()))
			}
		}
	}
}

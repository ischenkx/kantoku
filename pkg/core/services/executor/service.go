package executor

import (
	"context"
	"fmt"
	"github.com/ischenkx/kantoku/pkg/common/broker"
	codec "github.com/ischenkx/kantoku/pkg/common/data/codec"
	"github.com/ischenkx/kantoku/pkg/common/service"
	"github.com/ischenkx/kantoku/pkg/core/event"
	"github.com/ischenkx/kantoku/pkg/core/system"
	"github.com/ischenkx/kantoku/pkg/core/system/events"
	"golang.org/x/sync/errgroup"
	"log/slog"
)

const QueueName = "executor"

type Service struct {
	System      system.AbstractSystem
	ResultCodec codec.Codec[Result, []byte]
	Executor    Executor

	service.Core
}

func (srvc *Service) Run(ctx context.Context) error {
	g, ctx := errgroup.WithContext(ctx)

	executionService := &executionController{
		System:      srvc.System,
		Executor:    srvc.Executor,
		ResultCodec: srvc.ResultCodec,
		Service:     srvc.Core,
	}

	g.Go(func() error {
		readyTaskEvents, err := srvc.System.Events().Consume(ctx, broker.TopicsInfo{
			Group:  QueueName,
			Topics: []string{events.OnTask.Ready},
		})
		if err != nil {
			return fmt.Errorf("failed to read events: %w", err)
		}

		broker.Processor[event.Event]{
			Handler: func(ctx context.Context, ev event.Event) error {
				taskId := string(ev.Data)

				if err := executionService.processReadyTask(ctx, taskId); err != nil {
					return err
				}

				return nil
			},
			ErrorHandler: func(ctx context.Context, ev event.Event, err error) {
				taskId := string(ev.Data)

				srvc.Logger().
					Error("failed to process a ready task",
						slog.String("id", taskId),
						slog.String("error", err.Error()))
			},
		}.Process(ctx, readyTaskEvents)

		return nil
	})

	g.Go(func() error {
		if err := executionService.start(ctx); err != nil {
			return fmt.Errorf("failed to start the controller: %w", err)
		}
		return nil
	})

	if err := g.Wait(); err != nil {
		return err
	}

	return nil
}

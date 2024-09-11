package executor

import (
	"context"
	"fmt"
	codec "github.com/ischenkx/kantoku/pkg/common/data/codec"
	"github.com/ischenkx/kantoku/pkg/common/service"
	"github.com/ischenkx/kantoku/pkg/common/transport/broker"
	"github.com/ischenkx/kantoku/pkg/core"
	"golang.org/x/sync/errgroup"
	"log/slog"
)

const QueueName = "executor"

type Service struct {
	System      core.AbstractSystem
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

	readyTaskEvents, err := srvc.System.Events().Consume(ctx,
		[]string{core.OnTask.Ready},
		broker.ConsumerSettings{Group: QueueName},
	)
	if err != nil {
		return fmt.Errorf("failed to read events: %w", err)
	}

	for i := 0; i < 100; i++ {
		i := i
		g.Go(func() error {
			srvc.Logger().Info("starting a processor",
				"worker", i+1)
			broker.Processor[core.Event]{
				Handler: func(ctx context.Context, ev core.Event) error {
					taskId := string(ev.Data)

					srvc.Logger().Info("received a task", "task_id", taskId)

					if err := executionService.processReadyTask(ctx, taskId); err != nil {
						return err
					}

					return nil
				},
				ErrorHandler: func(ctx context.Context, ev core.Event, err error) {
					taskId := string(ev.Data)

					srvc.Logger().
						Error("failed to process a ready task",
							slog.String("id", taskId),
							slog.String("error", err.Error()))
				},
			}.Process(ctx, readyTaskEvents)

			return nil
		})
	}

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

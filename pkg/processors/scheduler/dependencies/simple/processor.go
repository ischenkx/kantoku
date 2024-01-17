package simple

import (
	"context"
	"fmt"
	"github.com/ischenkx/kantoku/pkg/common/data/deps"
	"github.com/ischenkx/kantoku/pkg/processors/scheduler/dependencies/simple/manager"
	"github.com/ischenkx/kantoku/pkg/system"
	event2 "github.com/ischenkx/kantoku/pkg/system/kernel/event"
	"golang.org/x/sync/errgroup"
	"log/slog"
)

var QueueName = "dependencies.simple"

type Processor struct {
	system  system.AbstractSystem
	manager *manager.Manager
}

func NewProcessor(
	system system.AbstractSystem,
	dependencies deps.Manager,
	task2group manager.TaskToGroup,
	resolvers map[string]manager.Resolver,
) *Processor {
	return &Processor{
		system: system,
		manager: manager.New(
			system,
			dependencies,
			task2group,
			resolvers,
		),
	}
}

func (processor *Processor) Process(ctx context.Context) error {
	g, ctx := errgroup.WithContext(ctx)
	g.Go(func() error {
		slog.Info("processing incoming tasks...")
		return processor.processNewTasks(ctx)
	})

	g.Go(func() error {
		slog.Info("sending ready tasks...")
		return processor.processReadyTasks(ctx)
	})

	if err := g.Wait(); err != nil {
		return err
	}

	return nil
}

func (processor *Processor) processNewTasks(ctx context.Context) error {
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

			if err := processor.manager.Register(ctx, taskId); err != nil {
				slog.Info("failed to register task",
					slog.String("id", taskId),
					slog.String("error", err.Error()))
			}
		}
	}
}

func (processor *Processor) processReadyTasks(ctx context.Context) error {
	channel, err := processor.manager.Ready(ctx)
	if err != nil {
		return err
	}

	for {
		select {
		case <-ctx.Done():
			return nil
		case taskId := <-channel:
			slog.Info("processing a ready task",
				slog.String("id", taskId))
			err := processor.system.Events().Publish(ctx, event2.New(system.TaskReadyEvent, []byte(taskId)))
			if err != nil {
				slog.Info("failed to publish an event",
					slog.String("id", taskId),
					slog.String("event", system.TaskReadyEvent),
					slog.String("error", err.Error()))
			}
		}
	}
}

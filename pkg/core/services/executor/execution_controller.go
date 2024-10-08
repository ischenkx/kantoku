package executor

import (
	"context"
	"fmt"
	"github.com/ischenkx/kantoku/pkg/common/data/codec"
	"github.com/ischenkx/kantoku/pkg/common/service"
	"github.com/ischenkx/kantoku/pkg/common/transport/broker"
	"github.com/ischenkx/kantoku/pkg/core"
	"sync"
)

type process struct {
	cancel context.CancelFunc
}

type executionController struct {
	System      core.AbstractSystem
	Executor    Executor
	ResultCodec codec.Codec[Result, []byte]
	Service     service.Core

	runningProcesses map[string]process
	mu               sync.Mutex
}

func (controller *executionController) start(ctx context.Context) error {
	cancellationEvents, err := controller.System.Events().Consume(
		ctx,
		[]string{core.OnTask.Cancelled},
		broker.ConsumerSettings{Group: controller.Service.ID()},
	)
	if err != nil {
		return err
	}

	go controller.processCancellationEvents(ctx, cancellationEvents)

	return nil
}

func (controller *executionController) processCancellationEvents(ctx context.Context, cancellationEvents <-chan broker.Message[core.Event]) {
	broker.Processor[core.Event]{
		Handler: func(ctx context.Context, ev core.Event) error {
			taskId := string(ev.Data)
			controller.cancel(ctx, taskId)
			return nil
		},
	}.Process(ctx, cancellationEvents)
}

func (controller *executionController) processReadyTask(ctx context.Context, id string) error {
	err := controller.System.Events().Send(ctx, core.NewEvent(core.OnTask.Received, []byte(id)))
	if err != nil {
		return err
	}

	result := Result{TaskID: id, Status: OK}
	if err := controller.execute(ctx, id); err != nil {
		result.Data = []byte(err.Error())
		result.Status = Failed
		// TODO: may be remove
		controller.Service.Logger().Warn("execution error:", err)
	}

	encodedResult, err := controller.ResultCodec.Encode(result)
	if err != nil {
		return fmt.Errorf("failed to encode the result: %w", err)
	}

	err = controller.System.Events().Send(ctx, core.NewEvent(core.OnTask.Finished, encodedResult))
	if err != nil {
		return err
	}

	return nil
}

func (controller *executionController) execute(ctx context.Context, id string) error {
	localContext, cancel := context.WithCancel(ctx)
	controller.createProcess(id, cancel)
	defer controller.deleteProcess(id)

	t, err := controller.System.Task(ctx, id)
	if err != nil {
		return fmt.Errorf("failed to load task: %w", err)
	}

	if err := controller.validateReadyTask(localContext, t); err != nil {
		return fmt.Errorf("failed to validate a task: %w", err)
	}

	err = controller.Executor.Execute(localContext, controller.System, t)
	if err != nil {
		return err
	}

	return nil
}

func (controller *executionController) validateReadyTask(ctx context.Context, t core.Task) error {
	if rawStatus, ok := t.Info["status"]; ok {
		if value, ok := rawStatus.(string); ok && value == core.TaskStatuses.Cancelled {
			return fmt.Errorf("task canceled")
		}
	}

	return nil
}

func (controller *executionController) cancel(_ context.Context, id string) {
	controller.mu.Lock()
	defer controller.mu.Unlock()

	if proc, ok := controller.runningProcesses[id]; ok {
		proc.cancel()
	}
}

func (controller *executionController) createProcess(id string, cancel context.CancelFunc) {
	controller.mu.Lock()
	defer controller.mu.Unlock()

	if controller.runningProcesses == nil {
		controller.runningProcesses = map[string]process{}
	}

	controller.runningProcesses[id] = process{cancel: cancel}
}

func (controller *executionController) deleteProcess(id string) {
	controller.mu.Lock()
	defer controller.mu.Unlock()

	delete(controller.runningProcesses, id)
}

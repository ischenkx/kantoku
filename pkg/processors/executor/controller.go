package executor

import (
	"context"
	"fmt"
	"github.com/ischenkx/kantoku/pkg/common/data/codec"
	"github.com/ischenkx/kantoku/pkg/system"
	"github.com/ischenkx/kantoku/pkg/system/kernel/event"
	"github.com/ischenkx/kantoku/pkg/system/kernel/task"
	"sync"
)

type process struct {
	cancel context.CancelFunc
}

type controller struct {
	system           system.AbstractSystem
	executor         Executor
	resultCodec      codec.Codec[Result, []byte]
	runningProcesses map[string]process
	localQueue       string
	mu               sync.Mutex
}

func newExecutionController(
	system system.AbstractSystem,
	executor Executor,
	localQueue string,
	resultCodec codec.Codec[Result, []byte],
) *controller {
	return &controller{
		system:           system,
		executor:         executor,
		resultCodec:      resultCodec,
		runningProcesses: make(map[string]process),
		localQueue:       localQueue,
	}
}

func (controller *controller) start(ctx context.Context) error {
	cancellationEvents, err := controller.system.Events().Consume(ctx,
		event.Queue{
			Name:   controller.localQueue,
			Topics: []string{system.TaskCancelledEvent},
		})
	if err != nil {
		return err
	}

	go controller.processCancellationEvents(ctx, cancellationEvents)

	return nil
}

func (controller *controller) processCancellationEvents(ctx context.Context, cancellationEvents <-chan event.Event) {
	for {
		select {
		case <-ctx.Done():
			return
		case ev := <-cancellationEvents:
			taskId := string(ev.Data)
			controller.cancel(ctx, taskId)
		}
	}
}

func (controller *controller) processReadyTask(ctx context.Context, id string) error {
	err := controller.system.Events().Publish(ctx, event.New(system.TaskReceivedEvent, []byte(id)))
	if err != nil {
		return fmt.Errorf("failed to publish a 'task received' event: %w", err)
	}

	result := Result{TaskID: id, Status: OK}
	if err := controller.execute(ctx, id); err != nil {
		result.Data = []byte(err.Error())
		result.Status = Failed
	}

	encodedResult, err := controller.resultCodec.Encode(result)
	if err != nil {
		return fmt.Errorf("failed to encode the result: %w", err)
	}

	err = controller.system.Events().Publish(ctx, event.New(system.TaskFinishedEvent, encodedResult))
	if err != nil {
		return fmt.Errorf("failed to publish a 'task_finished' event: %w", err)
	}

	return nil
}

func (controller *controller) execute(ctx context.Context, id string) error {
	localContext, cancel := context.WithCancel(ctx)
	controller.createProcess(id, cancel)
	defer controller.deleteProcess(id)

	task := controller.system.Task(id)

	if err := controller.validateReadyTask(localContext, task); err != nil {
		return fmt.Errorf("failed to validate a task: %w", err)
	}

	rawTask, err := task.Raw(localContext)
	if err != nil {
		return fmt.Errorf("failed to load a task: %w", err)
	}

	err = controller.executor.Execute(localContext, controller.system, rawTask)
	if err != nil {
		return err
	}

	return nil
}

func (controller *controller) validateReadyTask(ctx context.Context, t *system.Task) error {
	info, err := t.Info(ctx)
	if err != nil {
		return fmt.Errorf("failed to load task info: %w", err)
	}

	if rawStatus, ok := info["status"]; ok {
		if value, ok := rawStatus.(string); ok && value == task.CancelledStatus {
			return fmt.Errorf("task canceled")
		}
	}

	return nil
}

func (controller *controller) cancel(_ context.Context, id string) {
	controller.mu.Lock()
	defer controller.mu.Unlock()

	if proc, ok := controller.runningProcesses[id]; ok {
		proc.cancel()
	}
}

func (controller *controller) createProcess(id string, cancel context.CancelFunc) {
	controller.mu.Lock()
	defer controller.mu.Unlock()

	controller.runningProcesses[id] = process{cancel: cancel}
}

func (controller *controller) deleteProcess(id string) {
	controller.mu.Lock()
	defer controller.mu.Unlock()

	delete(controller.runningProcesses, id)
}

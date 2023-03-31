package kantoku

import (
	"context"
	"github.com/google/uuid"
	"kantoku/common/data/kv"
	"kantoku/core/event"
	"kantoku/core/task"
	"kantoku/framework/argument"
)

type Kantoku struct {
	plugins         []Plugin
	argumentManager *argument.Manager
	tasks           kv.Database[StoredTask]
	scheduler       *task.Scheduler[ScheduledTask]
	properties      *Properties
	cells           Cells
	events          event.Bus
}

func New(
	argumentManager *argument.Manager,
	taskDB kv.Database[StoredTask],
	taskQueue task.SchedulerInputs[ScheduledTask],
	events event.Bus,
	cells Cells,
) *Kantoku {
	return &Kantoku{
		argumentManager: argumentManager,
		tasks:           taskDB,
		events:          events,
		scheduler:       task.NewScheduler(taskQueue, events),
		cells:           cells,
	}
}

func (kantoku *Kantoku) Spawn(ctx_ context.Context, spec Spec) (result Result, err error) {
	ctx := kantoku.makeContext(ctx_)
	result.Log = ctx.Log
	defer ctx.finalize()

	ctx.Task = &TaskInstance{
		id:        "",
		typ:       spec.Type,
		arguments: append([]any(nil), spec.Arguments...),
	}

	for index := 0; index < ctx.Task.CountArgs(); index++ {
		arg, _ := ctx.Task.Arg(index)
		if initializer, ok := arg.(ArgumentInitializer); ok {
			arg, err := initializer.Initialize(ctx)
			if err != nil {
				return result, err
			}
			ctx.Task.arguments[index] = arg
		}
	}

	for _, plugin := range kantoku.plugins {
		if err := plugin.BeforeInitialized(ctx); err != nil {
			return result, err
		}
	}

	storedTask := StoredTask{
		ID:   uuid.New().String(),
		Type: ctx.Task.Type(),
	}

	for index := 0; index < ctx.Task.CountArgs(); index++ {
		arg, _ := ctx.Task.Arg(index)
		formattedArgument, err := kantoku.argumentManager.Encode(ctx, arg)
		if err != nil {
			return result, err
		}
		storedTask.Arguments = append(storedTask.Arguments, formattedArgument)
	}

	if _, err := kantoku.tasks.Set(ctx, storedTask.ID, storedTask); err != nil {
		return result, err
	}
	result.Task = storedTask.ID
	ctx.Task.id = storedTask.ID

	for _, plugin := range kantoku.plugins {
		plugin.AfterInitialized(ctx)
	}

	for _, plugin := range kantoku.plugins {
		plugin.BeforeScheduled(ctx)
	}

	if err := kantoku.scheduler.Schedule(ctx, ScheduledTask{id: ctx.Task.ID()}); err != nil {
		return result, err
	}

	for _, plugin := range kantoku.plugins {
		plugin.AfterScheduled(ctx)
	}

	ctx.Log.Spawned = append(result.Log.Spawned, storedTask.ID)

	return result, nil
}

func (kantoku *Kantoku) Register(plugin Plugin) {
	plugin.Initialize(kantoku)
	kantoku.plugins = append(kantoku.plugins, plugin)
}

func (kantoku *Kantoku) Task(id string) TaskView {
	return TaskView{
		kantoku: kantoku,
		id:      id,
	}
}

func (kantoku *Kantoku) Events() event.Bus {
	return kantoku.events
}

func (kantoku *Kantoku) Cells() Cells {
	return kantoku.cells
}

func (kantoku *Kantoku) Props() *Properties {
	return kantoku.properties
}

func (kantoku *Kantoku) makeContext(ctx_ context.Context) *Context {
	ctx := NewContext(ctx_)
	ctx.parent = ctx_
	return ctx
}

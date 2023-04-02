package kantoku

import (
	"context"
	"github.com/google/uuid"
	"kantoku/common/data/kv"
	"kantoku/core/event"
)

type Builder struct {
	Tasks     kv.Database[StoredTask]
	Scheduler Scheduler
	Cells     Cells
	Events    event.Bus
}

func (builder Builder) Build() *Kantoku {
	return &Kantoku{
		tasks:      builder.Tasks,
		scheduler:  builder.Scheduler,
		cells:      builder.Cells,
		events:     builder.Events,
		properties: NewProperties(),
		plugins:    nil,
	}
}

type Kantoku struct {
	tasks      kv.Database[StoredTask]
	scheduler  Scheduler
	cells      Cells
	events     event.Bus
	properties *Properties
	plugins    []Plugin
}

func (kantoku *Kantoku) Spawn(ctx_ context.Context, spec Spec) (result Result, err error) {
	ctx := kantoku.makeContext(ctx_)
	result.Log = ctx.Log
	defer ctx.finalize()

	ctx.Task = &TaskInstance{
		id:   "",
		typ:  spec.Type,
		data: spec.Data,
	}

	if initializer, ok := ctx.Task.Data().(Initializeable); ok {
		data, err := initializer.Initialize(ctx)
		if err != nil {
			return result, err
		}
		ctx.Task.data = data
	}

	for _, plugin := range kantoku.plugins {
		if err := plugin.BeforeInitialized(ctx); err != nil {
			return result, err
		}
	}

	storedTask := StoredTask{
		Id:   uuid.New().String(),
		Type: ctx.Task.Type(),
		Data: ctx.Task.Data(),
	}

	if _, err := kantoku.tasks.Set(ctx, storedTask.Id, storedTask); err != nil {
		return result, err
	}
	result.Task = storedTask.Id
	ctx.Task.id = storedTask.Id

	for _, plugin := range kantoku.plugins {
		plugin.AfterInitialized(ctx)
	}

	for _, plugin := range kantoku.plugins {
		plugin.BeforeScheduled(ctx)
	}

	if err := kantoku.scheduler.Schedule(ctx); err != nil {
		return result, err
	}

	for _, plugin := range kantoku.plugins {
		plugin.AfterScheduled(ctx)
	}

	ctx.Log.Spawned = append(result.Log.Spawned, storedTask.Id)

	return result, nil
}

func (kantoku *Kantoku) Register(plugin Plugin) {
	plugin.Initialize(kantoku)
	kantoku.plugins = append(kantoku.plugins, plugin)
}

func (kantoku *Kantoku) Task(id string) *TaskView {
	return &TaskView{
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

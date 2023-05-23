package kantoku

import (
	"context"
	"github.com/google/uuid"
	"kantoku/common/data/kv"
	"kantoku/platform"
)

type Builder struct {
	Tasks   kv.Database[string, TaskInstance]
	Inputs  platform.Inputs[TaskInstance]
	Outputs platform.Outputs
	Events  platform.Broker
}

func (builder Builder) Build() *Kantoku {
	return &Kantoku{
		tasks:      builder.Tasks,
		outputs:    builder.Outputs,
		inputs:     builder.Inputs,
		events:     builder.Events,
		properties: NewProperties(),
		plugins:    nil,
	}
}

type Kantoku struct {
	tasks      kv.Database[string, TaskInstance]
	outputs    platform.Outputs
	inputs     platform.Inputs[TaskInstance]
	events     platform.Broker
	properties *Properties
	plugins    []Plugin
}

func (kantoku *Kantoku) Spawn(ctx_ context.Context, spec Spec) (result Result, err error) {
	ctx := kantoku.makeContext(ctx_)
	result.Log = ctx.Log
	defer ctx.finalize()

	ctx.Task = &TaskInstance{
		id:  "",
		typ: spec.Type,
	}

	for _, option := range spec.Options {
		if err := option(ctx); err != nil {
			return result, err
		}
	}

	for _, plugin := range kantoku.plugins {
		if p, ok := plugin.(BeforeInitializedPlugin); ok {
			if err := p.BeforeInitialized(ctx); err != nil {
				return result, err
			}
		}
	}

	if err := kantoku.tasks.Set(ctx, ctx.Task.ID, ctx.Task); err != nil {
		return result, err
	}
	result.Task = ctx.Task.ID

	for _, plugin := range kantoku.plugins {
		if p, ok := plugin.(AfterInitializedPlugin); ok {
			p.AfterInitialized(ctx)
		}
	}

	for _, plugin := range kantoku.plugins {
		if p, ok := plugin.(BeforeScheduledPlugin); ok {
			if err := p.BeforeScheduled(ctx); err != nil {
				return result, err
			}
		}
	}

	if err := kantoku.inputs.Write(ctx, ctx.Task); err != nil {
		return result, err
	}

	for _, plugin := range kantoku.plugins {
		if p, ok := plugin.(AfterScheduledPlugin); ok {
			p.AfterScheduled(ctx)
		}
	}

	ctx.Log.Spawned = append(result.Log.Spawned, ctx.Task.ID)

	return result, nil
}

func (kantoku *Kantoku) Register(plugin Plugin) {
	if p, ok := plugin.(InitializePlugin); ok {
		p.Initialize(*kantoku)
	}
	kantoku.plugins = append(kantoku.plugins, plugin)
}

func (kantoku *Kantoku) Task(id string) *View {
	return &View{
		kantoku: kantoku,
		id:      id,
	}
}

func (kantoku *Kantoku) Events() event.Bus {
	return kantoku.events
}

func (kantoku *Kantoku) Props() *Properties {
	return kantoku.properties
}

func (kantoku *Kantoku) makeContext(ctx_ context.Context) *Context {
	ctx := NewContext(ctx_)
	ctx.parent = ctx_
	return ctx
}

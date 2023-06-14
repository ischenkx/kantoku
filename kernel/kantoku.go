package kernel

import (
	"context"
	"github.com/google/uuid"
	"kantoku/kernel/platform"
)

type Kernel struct {
	platform platform.Platform[Task]
	plugins  []Plugin
}

func New(platform platform.Platform[Task]) *Kernel {
	return &Kernel{
		platform: platform,
	}
}

func (kantoku *Kernel) Spawn(ctx_ context.Context, spec Spec) (result Result, err error) {
	ctx := kantoku.makeContext(ctx_)
	result.Log = ctx.Log
	defer ctx.finalize()

	ctx.Task = Task{
		Id:   uuid.New().String(),
		Type: spec.Type,
		Data: spec.Data,
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

	if err := kantoku.platform.DB().Set(ctx, ctx.Task.ID(), ctx.Task); err != nil {
		return result, err
	}
	result.Task = ctx.Task.ID()

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

	if err := kantoku.platform.Inputs().Write(ctx, ctx.Task.ID()); err != nil {
		return result, err
	}

	for _, plugin := range kantoku.plugins {
		if p, ok := plugin.(AfterScheduledPlugin); ok {
			p.AfterScheduled(ctx)
		}
	}

	ctx.Log.Spawned = append(result.Log.Spawned, ctx.Task.ID())

	return result, nil
}

func (kantoku *Kernel) Register(plugin Plugin) {
	if p, ok := plugin.(InitializePlugin); ok {
		p.Initialize(*kantoku)
	}
	kantoku.plugins = append(kantoku.plugins, plugin)
}

func (kantoku *Kernel) Task(id string) *View {
	return &View{
		kantoku: kantoku,
		id:      id,
	}
}

func (kantoku *Kernel) Broker() platform.Broker {
	return kantoku.platform.Broker()
}
func (kantoku *Kernel) Outputs() platform.Outputs {
	return kantoku.platform.Outputs()
}

func (kantoku *Kernel) makeContext(ctx_ context.Context) *Context {
	ctx := NewContext(ctx_)
	ctx.parent = ctx_
	return ctx
}

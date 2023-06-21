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
func (kernel *Kernel) Spawn(ctx_ context.Context, spec Spec) (result Result, err error) {
	ctx := kernel.makeContext(ctx_)
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

	for _, plugin := range kernel.plugins {
		if p, ok := plugin.(BeforeInitializedPlugin); ok {
			if err := p.BeforeInitialized(ctx); err != nil {
				return result, err
			}
		}
	}

	if err := kernel.platform.DB().Set(ctx, ctx.Task.ID(), ctx.Task); err != nil {
		return result, err
	}
	result.Task = ctx.Task.ID()

	for _, plugin := range kernel.plugins {
		if p, ok := plugin.(AfterInitializedPlugin); ok {
			p.AfterInitialized(ctx)
		}
	}

	for _, plugin := range kernel.plugins {
		if p, ok := plugin.(BeforeScheduledPlugin); ok {
			if err := p.BeforeScheduled(ctx); err != nil {
				return result, err
			}
		}
	}

	if err := kernel.platform.Inputs().Write(ctx, ctx.Task.ID()); err != nil {
		return result, err
	}

	for _, plugin := range kernel.plugins {
		if p, ok := plugin.(AfterScheduledPlugin); ok {
			p.AfterScheduled(ctx)
		}
	}

	ctx.Log.Spawned = append(result.Log.Spawned, ctx.Task.ID())

	return result, nil
}
func (kernel *Kernel) Register(plugin Plugin) error {
	for _, plugin := range kernel.plugins {
		if p, ok := plugin.(BeforePluginInitPlugin); ok {
			if err := p.BeforePluginInit(kernel, plugin); err != nil {
				return err
			}
		}
	}
	if p, ok := plugin.(InitializePlugin); ok {
		p.Initialize(*kernel)
	}
	kernel.plugins = append(kernel.plugins, plugin)

	for _, plugin := range kernel.plugins {
		if p, ok := plugin.(AfterPluginInitPlugin); ok {
			p.AfterPluginInit(kernel, plugin)
		}
	}

	return nil
}
func (kernel *Kernel) Task(id string) *View {
	return &View{
		kernel: kernel,
		id:     id,
	}
}
func (kernel *Kernel) Broker() platform.Broker {
	return kernel.platform.Broker()
}
func (kernel *Kernel) Outputs() platform.Outputs {
	return kernel.platform.Outputs()
}
func (kernel *Kernel) Plugins() []Plugin {
	return kernel.plugins
}
func (kernel *Kernel) makeContext(ctx_ context.Context) *Context {
	ctx := NewContext(ctx_)
	ctx.parent = ctx_
	return ctx
}

package job

import (
	"context"
	"github.com/google/uuid"
	"kantoku/common/data/kv"
	"kantoku/common/data/pool"
)

type Inputs pool.Pool[string]
type Outputs kv.Database[string, Result]
type DB kv.Database[string, Job]

type Manager struct {
	inputs  Inputs
	outputs Outputs
	db      DB
	plugins []Plugin
}

func NewManager(inputs Inputs, outputs Outputs, db DB) *Manager {
	return &Manager{
		inputs:  inputs,
		outputs: outputs,
		db:      db,
	}
}

func (manager *Manager) Spawn(ctx_ context.Context, spec Spec) (result SpawnResult, err error) {
	ctx := manager.makeContext(ctx_)
	result.Log = ctx.Log
	defer ctx.finalize()

	ctx.Task = Job{
		Id:   uuid.New().String(),
		Type: spec.Type,
		Data: spec.Data,
	}

	for _, option := range spec.Options {
		if err := option(ctx); err != nil {
			return result, err
		}
	}

	for _, plugin := range manager.plugins {
		if p, ok := plugin.(BeforeInitializedPlugin); ok {
			if err := p.BeforeInitialized(ctx); err != nil {
				return result, err
			}
		}
	}

	if err := manager.db.Set(ctx, ctx.Task.ID(), ctx.Task); err != nil {
		return result, err
	}
	result.Task = ctx.Task.ID()

	for _, plugin := range manager.plugins {
		if p, ok := plugin.(AfterInitializedPlugin); ok {
			p.AfterInitialized(ctx)
		}
	}

	for _, plugin := range manager.plugins {
		if p, ok := plugin.(BeforeScheduledPlugin); ok {
			if err := p.BeforeScheduled(ctx); err != nil {
				return result, err
			}
		}
	}

	if err := manager.Inputs().Write(ctx, ctx.Task.ID()); err != nil {
		return result, err
	}

	for _, plugin := range manager.plugins {
		if p, ok := plugin.(AfterScheduledPlugin); ok {
			p.AfterScheduled(ctx)
		}
	}

	ctx.Log.Spawned = append(result.Log.Spawned, ctx.Task.ID())

	return result, nil
}

func (manager *Manager) Use(plugin Plugin) error {
	for _, plugin := range manager.plugins {
		if p, ok := plugin.(BeforePluginInitPlugin); ok {
			if err := p.BeforePluginInit(manager, plugin); err != nil {
				return err
			}
		}
	}
	if p, ok := plugin.(InitializePlugin); ok {
		p.Initialize(manager)
	}
	manager.plugins = append(manager.plugins, plugin)

	for _, plugin := range manager.plugins {
		if p, ok := plugin.(AfterPluginInitPlugin); ok {
			p.AfterPluginInit(manager, plugin)
		}
	}

	return nil
}

func (manager *Manager) Task(id string) *View {
	return &View{
		kernel: manager,
		id:     id,
	}
}

func (manager *Manager) Inputs() Inputs {
	return manager.inputs
}

func (manager *Manager) Outputs() Outputs {
	return manager.outputs
}

func (manager *Manager) Tasks() kv.Getter[string, Job] {
	return manager.db
}

func (manager *Manager) Plugins() []Plugin {
	return manager.plugins
}

func (manager *Manager) makeContext(ctx_ context.Context) *Context {
	ctx := NewContext(ctx_)
	ctx.parent = ctx_
	return ctx
}

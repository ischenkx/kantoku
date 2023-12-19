package alpha

import (
	"context"
	"fmt"
	"kantoku/common/data/identifier"
	"kantoku/common/data/pool"
	"kantoku/common/data/record"
	"kantoku/core/daemon"
	"kantoku/core/util"
)

type Manager struct {
	pool    Pool
	results ResultStorage
	storage Storage
	ids     identifier.Generator
	info    record.Storage
	runner  Runner
	plugins []Plugin
}

func (manager *Manager) Info() record.Storage {
	return manager.info
}

func (manager *Manager) Get(id string) Alpha {
	return Alpha{id: id}
}

func (manager *Manager) New(ctx context.Context, data []byte) (alpha Alpha, err error) {
	id, err := manager.ids.New(ctx)
	if err != nil {
		return alpha, fmt.Errorf("failed to generate an id: %s", err)
	}

	if err := manager.storage.Set(ctx, id, data); err != nil {
		return alpha, fmt.Errorf("failed to initialize the process: %s", err)
	}

	for _, plugin := range util.FilterPlugins[OnNewPlugin](manager.Plugins()) {
		plugin.OnNew(ctx, manager.Get(id))
	}

	return Alpha{id: id}, nil
}

func (manager *Manager) Run(ctx context.Context, id string) error {
	for _, plugin := range util.FilterPlugins[OnBeforeRunPlugin](manager.Plugins()) {
		if err := plugin.OnBeforeRun(ctx, manager.Get(id)); err != nil {
			return fmt.Errorf("plugin failed: %s", err)
		}
	}

	if err := manager.pool.Write(ctx, id); err != nil {
		return fmt.Errorf("failed to put alpha in a pool: %s", err)
	}

	for _, plugin := range util.FilterPlugins[OnRunPlugin](manager.Plugins()) {
		plugin.OnRun(ctx, manager.Get(id))
	}

	return nil
}

func (manager *Manager) Daemons() []daemon.Daemon {
	return []daemon.Daemon{
		{
			Settings: daemon.Settings{Scalable: true},
			Func:     manager.Process,
			Type:     "ALPHA_EXECUTOR",
		},
	}
}

func (manager *Manager) Process(ctx context.Context) error {
	return pool.AutoCommit[string](ctx, manager.pool, func(ctx context.Context, id string) error {
		for _, plugin := range util.FilterPlugins[OnReceivePlugin](manager.Plugins()) {
			if err := plugin.OnReceive(ctx, manager.Get(id)); err != nil {
				return fmt.Errorf("on-receive plugin failed: %s", err)
			}
		}

		data, err := manager.runner.Run(ctx, manager.Get(id))

		result := Result{
			Data:   data,
			Status: OK,
		}

		if err != nil {
			result.Data = []byte(err.Error())
			result.Status = FAILURE
		}

		for _, plugin := range util.FilterPlugins[OnExecutedPlugin](manager.Plugins()) {
			plugin.OnExecuted(ctx, manager.Get(id), result)
		}

		if err := manager.results.Set(ctx, id, result); err != nil {
			for _, plugin := range util.FilterPlugins[OnResultSaveFailurePlugin](manager.Plugins()) {
				plugin.OnResultSaveFailure(ctx, manager.Get(id), result, err)
			}
			return err
		}

		return nil
	})
}

func (manager *Manager) Plugins() []Plugin {
	return manager.plugins
}

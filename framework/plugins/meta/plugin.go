package meta

import (
	"fmt"
	"kantoku/kernel"
)

type PluginData map[string]any

type Plugin struct {
	manager *Manager
}

func NewPlugin(manager *Manager) Plugin {
	return Plugin{manager: manager}
}

func (p Plugin) BeforeInitialized(ctx *kernel.Context) error {
	data := kernel.GetPluginData(ctx).GetWithDefault("meta", PluginData{}).(PluginData)

	if len(data) == 0 {
		return nil
	}

	meta, err := p.manager.Get(ctx, ctx.Task.ID())
	if err != nil {
		return fmt.Errorf("failed to get meta information: %s", err)
	}

	for key, value := range data {
		if err := meta.Get(key).Set(ctx, value); err != nil {
			return fmt.Errorf("failed to initialized a meta field: %s", err)
		}
	}

	return nil
}

func WithEntry(key string, value any) kernel.Option {
	return func(ctx *kernel.Context) error {
		data := kernel.GetPluginData(ctx).GetWithDefault("meta", PluginData{}).(PluginData)
		data[key] = value
		return nil
	}
}

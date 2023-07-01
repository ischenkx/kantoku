package info

import (
	"fmt"
	"kantoku/kernel"
)

type PluginData map[string]any

type Plugin struct {
	storage *Storage
}

func NewPlugin(storage *Storage) Plugin {
	return Plugin{storage: storage}
}

func (p Plugin) BeforeInitialized(ctx *kernel.Context) error {
	data := kernel.GetPluginData(ctx).GetWithDefault("info", PluginData{}).(PluginData)
	if len(data) == 0 {
		return nil
	}

	info := p.storage.Get(ctx.Task.ID())

	err := info.Set(ctx, Dict(data).AsEntries()...)
	if err != nil {
		return fmt.Errorf("failed to set info properties: %s", err)
	}

	return nil
}

func WithEntry(key string, value any) kernel.Option {
	return func(ctx *kernel.Context) error {
		data := kernel.GetPluginData(ctx).GetWithDefault("info", PluginData{}).(PluginData)
		data[key] = value
		return nil
	}
}

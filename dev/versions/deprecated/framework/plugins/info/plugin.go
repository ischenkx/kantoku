package info

import (
	"fmt"
	"kantoku/framework/job"
)

type PluginData map[string]any

type Plugin struct {
	storage *Storage
}

func NewPlugin(storage *Storage) Plugin {
	return Plugin{storage: storage}
}

func (p Plugin) BeforeInitialized(ctx *job.Context) error {
	data := job.GetPluginData(ctx).GetWithDefault("info", PluginData{}).(PluginData)
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

func WithEntry(key string, value any) job.Option {
	return func(ctx *job.Context) error {
		data := job.GetPluginData(ctx).GetWithDefault("info", PluginData{}).(PluginData)
		data[key] = value
		return nil
	}
}

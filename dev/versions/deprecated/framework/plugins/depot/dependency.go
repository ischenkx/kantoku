package depot

import (
	job2 "kantoku/framework/job"
)

func Dependency(id string) job2.Option {
	return func(ctx *job2.Context) error {
		data := ctx.Data().GetOrSet("dependencies", func() any { return &PluginData{} }).(*PluginData)
		data.Dependencies = append(data.Dependencies, id)
		return nil
	}
}

package depot

import (
	"kantoku/kernel"
)

func Dependency(id string) kernel.Option {
	return func(ctx *kernel.Context) error {
		data := ctx.Data().GetOrSet("dependencies", func() any { return &PluginData{} }).(*PluginData)
		data.Dependencies = append(data.Dependencies, id)
		return nil
	}
}

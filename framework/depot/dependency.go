package depot

import "kantoku"

func Dependency(id string) kantoku.Option {
	return func(ctx *kantoku.Context) error {
		data := ctx.Data().GetOrSet("dependencies", func() any { return &PluginData{} }).(*PluginData)
		data.Dependencies = append(data.Dependencies, id)
		return nil
	}
}

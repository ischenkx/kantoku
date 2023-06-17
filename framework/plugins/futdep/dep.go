package futdep

import (
	"kantoku/framework/future"
	"kantoku/kernel"
)

func Dep(id future.ID) kernel.Option {
	return func(ctx *kernel.Context) error {
		data := ctx.Data().GetOrSet("futdep", func() any { return &PluginData{} }).(*PluginData)
		data.Futures = append(data.Futures, id)

		return nil
	}
}

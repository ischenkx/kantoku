package futdep

import (
	"kantoku/common/data/future"
	job2 "kantoku/framework/job"
)

func Dep(id future.ID) job2.Option {
	return func(ctx *job2.Context) error {
		data := ctx.Data().GetOrSet("futdep", func() any { return &PluginData{} }).(*PluginData)
		data.Futures = append(data.Futures, id)

		return nil
	}
}

package delay

import (
	"kantoku"
	"time"
)

func Delay(duration time.Duration) kantoku.Option {
	return func(ctx *kantoku.Context) error {
		data := ctx.Data().GetOrSet("delay", func() any { return &PluginData{} }).(*PluginData)
		when := time.Now().Add(duration)
		data.When = &when
		return nil
	}
}

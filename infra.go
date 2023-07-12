package kantoku

import (
	"context"
	"github.com/samber/lo"
	demons2 "kantoku/framework/infra/demon"
	"kantoku/framework/utils/demons"
	"kantoku/kernel"
)

type Infra struct {
	kantoku *Kantoku
}

func (infra Infra) Demons(ctx context.Context) []demons2.Demon {
	providers := lo.Map(infra.kantoku.Kernel().Plugins(), func(plugin kernel.Plugin, _ int) demons2.Provider {
		return demons.TryProvider(plugin)
	})

	return demons.Multi(providers).Demons(ctx)
}

func (infra Infra) Deploy(ctx context.Context) error {
	return demons.Deploy(ctx, demons.List(infra.Demons(ctx)), infra.kantoku.deployer)
}

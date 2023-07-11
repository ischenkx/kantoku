package demons

import (
	"context"
	"kantoku/framework/infra/demon"
)

func Deploy(ctx context.Context, provider demon.Provider, orchestrator demon.Deployer) error {
	return demon.NewManager(orchestrator).Register(provider.Demons(ctx)...).Run(ctx)
}

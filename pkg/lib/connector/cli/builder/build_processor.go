package builder

import (
	"context"
	"fmt"
	"github.com/ischenkx/kantoku/pkg/common/service"
	"github.com/ischenkx/kantoku/pkg/core/system"
	"github.com/ischenkx/kantoku/pkg/lib/connector/cli/config"
)

func (builder *Builder) BuildProcessor(ctx context.Context, sys system.AbstractSystem, cfg config.DynamicConfig) (service.Service, []service.Middleware, error) {
	return nil, nil, fmt.Errorf("not supported")
}

package builder

import (
	"context"
	"github.com/ischenkx/kantoku/pkg/common/data/uid"
	"github.com/ischenkx/kantoku/pkg/common/service"
	"github.com/ischenkx/kantoku/pkg/lib/connector/cli/config"
	"github.com/ischenkx/kantoku/pkg/lib/connector/cli/errx"
	"os"
)

type ServiceData[T any] struct {
	Info T `mapstructure:"$"`
}

func (builder *Builder) BuildServiceCore(ctx context.Context, defaultName string, cfg config.DynamicConfig) (service.Core, error) {
	var data ServiceData[struct {
		Name string
		ID   string
	}]
	if err := cfg.Bind(&data); err != nil {
		return service.Core{}, errx.FailedToBind(err)
	}

	if data.Info.Name == "" {
		data.Info.Name = defaultName
	}

	if data.Info.ID == "" {
		data.Info.ID = uid.Generate()
	}

	core := service.NewCore(
		data.Info.Name,
		data.Info.ID,
		newLogger(os.Stdout).With("service", data.Info.Name),
	)

	return core, nil
}

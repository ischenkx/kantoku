package exec

import (
	"context"
	"errors"
	"kantoku/backend/executor"
	demons2 "kantoku/common/util/demons"
	"kantoku/framework/infra"
	"kantoku/framework/job"
)

func New(runner executor.Runner) *Plugin {
	return &Plugin{
		runner: runner,
	}
}

type Plugin struct {
	runner executor.Runner
	kernel *job.Manager
}

func (p *Plugin) Initialize(kernel *job.Manager) {
	p.kernel = kernel
}

func (p *Plugin) Demons() []infra.Demon {
	return demons2.List{
		demons2.Functional("EXEC", func(ctx context.Context) error {
			if p.kernel == nil {
				return errors.New("no kernel provided")
			}

			var plugins []executor.Plugin

			for _, plugin := range p.kernel.Plugins() {
				if provider, ok := plugin.(PluginProvider); ok {
					plugins = append(plugins, provider.ExecutorPlugins()...)
				}
			}

			return executor.New(p.kernel, p.runner).Use(plugins...).Run(ctx)
		}),
	}
}

type PluginProvider interface {
	ExecutorPlugins() []executor.Plugin
}

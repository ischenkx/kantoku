package deprecated

import (
	"github.com/samber/lo"
	"kantoku/common/codec"
	"kantoku/common/data/future"
	"kantoku/framework/infra"
	"kantoku/framework/job"
	"kantoku/framework/plugins/info"
)

type Kantoku struct {
	info                 *info.Storage
	futures              *future.Manager
	jobs                 *job.Manager
	dependencies         DependencyManager
	parametrizationCodec codec.Codec[Parametrization, []byte]
	settings             Settings
}

func (kantoku *Kantoku) Tasks() TaskManager {
	return TaskManager{kantoku}
}

func (kantoku *Kantoku) Dependencies() DependencyManager {
	return kantoku.dependencies
}

func (kantoku *Kantoku) Futures() *future.Manager {
	return kantoku.futures
}

func (kantoku *Kantoku) Demons() []infra.Demon {
	return lo.Flatten[infra.Demon](lo.Map(
		kantoku.Plugins(), func(plugin job.Plugin, _ int) []infra.Demon {
			if provider, ok := plugin.(infra.Provider); ok {
				return provider.Demons()
			}
			return nil
		}))
}

func (kantoku *Kantoku) Plugins() []job.Plugin {
	return kantoku.jobs.Plugins()
}

func (kantoku *Kantoku) Jobs() *job.Manager {
	return kantoku.jobs
}

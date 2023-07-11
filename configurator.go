package kantoku

import (
	"kantoku/common/codec"
	"kantoku/common/data/bimap"
	"kantoku/common/data/kv"
	"kantoku/common/data/record"
	"kantoku/common/functions"
	"kantoku/framework/future"
	"kantoku/framework/infra/demon"
	"kantoku/framework/plugins/demonic"
	"kantoku/framework/plugins/depot"
	"kantoku/framework/plugins/depot/deps"
	"kantoku/framework/plugins/futdep"
	"kantoku/framework/plugins/info"
	"kantoku/framework/plugins/taskdep"
	"kantoku/kernel"
	"kantoku/kernel/platform"
)

type Configurator struct {
	info struct {
		storage  record.Storage
		settings info.Settings
	}

	deps deps.Deps

	depot struct {
		groupTaskBimap bimap.Bimap[string, string]
	}

	platform platform.Platform[kernel.Task]

	parametrizationCodec codec.Codec[Parametrization, []byte]

	futures *future.Manager

	futdep struct {
		future2dependency kv.Database[future.ID, string]
	}

	taskdep struct {
		task2dependency kv.Database[string, string]
		topic           string
	}

	executor struct {
		runner functions.Runner
	}

	settings Settings

	plugins []kernel.Plugin

	deployer demon.Deployer
}

func Configure() Configurator {
	return Configurator{}
}

func (configurator Configurator) Info(storage record.Storage, settings info.Settings) Configurator {
	configurator.info.storage = storage
	configurator.info.settings = settings
	return configurator
}

func (configurator Configurator) Deps(deps deps.Deps) Configurator {
	configurator.deps = deps
	return configurator
}

func (configurator Configurator) Depot(groupTaskBimap bimap.Bimap[string, string]) Configurator {
	configurator.depot.groupTaskBimap = groupTaskBimap
	return configurator
}

func (configurator Configurator) Platform(platform platform.Platform[kernel.Task]) Configurator {
	configurator.platform = platform
	return configurator
}

func (configurator Configurator) Parametrization(parametrizationCodec codec.Codec[Parametrization, []byte]) Configurator {
	configurator.parametrizationCodec = parametrizationCodec
	return configurator
}

func (configurator Configurator) Futures(futures *future.Manager) Configurator {
	configurator.futures = futures
	return configurator
}

func (configurator Configurator) FutDep(future2dependency kv.Database[future.ID, string]) Configurator {
	configurator.futdep.future2dependency = future2dependency
	return configurator
}

func (configurator Configurator) TaskDep(task2dependency kv.Database[string, string]) Configurator {
	configurator.taskdep.task2dependency = task2dependency
	return configurator
}

func (configurator Configurator) Settings(settings Settings) Configurator {
	configurator.settings = settings
	return configurator
}

func (configurator Configurator) Demons(demons ...demon.Demon) Configurator {
	return configurator.Plugins(demonic.New(demons...))
}

func (configurator Configurator) Deployer(deployer demon.Deployer) Configurator {
	configurator.deployer = deployer
	return configurator
}

func (configurator Configurator) Plugins(plugins ...kernel.Plugin) Configurator {
	configurator.plugins = append(configurator.plugins, plugins...)
	return configurator
}

func (configurator Configurator) Compile() (*Kantoku, error) {
	dep := depot.New(configurator.deps, configurator.depot.groupTaskBimap, configurator.platform.Inputs())

	p := platform.New[kernel.Task](
		configurator.platform.DB(),
		dep,
		configurator.platform.Outputs(),
		configurator.platform.Broker(),
	)
	kan := &Kantoku{
		depot:                dep,
		parametrizationCodec: configurator.parametrizationCodec,
		info:                 info.NewStorage(configurator.info.storage, configurator.info.settings),
		futures:              configurator.futures,
		futdep:               futdep.NewManager(configurator.deps, configurator.futdep.future2dependency),
		taskdep:              taskdep.NewManager(configurator.deps, configurator.taskdep.task2dependency),
		kernel:               kernel.New(p),
		deployer:             configurator.deployer,
		settings:             configurator.settings,
	}

	plugins := []kernel.Plugin{
		info.NewPlugin(kan.info),
		futdep.NewPlugin(kan.futdep, kan.futures),
		taskdep.NewPlugin(kan.taskdep, kan.kernel.Broker(), configurator.taskdep.topic),
		depot.NewPlugin(kan.depot),
	}
	plugins = append(plugins, configurator.plugins...)
	for _, plugin := range plugins {
		if err := kan.kernel.Register(plugin); err != nil {
			return nil, err
		}
	}

	return kan, nil
}

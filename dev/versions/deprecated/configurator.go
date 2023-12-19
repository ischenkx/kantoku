package deprecated

import (
	"kantoku/common/codec"
	"kantoku/common/data/bimap"
	"kantoku/common/data/deps"
	"kantoku/common/data/future"
	"kantoku/common/data/kv"
	"kantoku/common/data/record"
	job2 "kantoku/framework/job"
	"kantoku/framework/plugins/depot"
	"kantoku/framework/plugins/futdep"
	"kantoku/framework/plugins/info"
	"kantoku/framework/plugins/taskdep"
)

type Configurator struct {
	info struct {
		storage  record.Storage
		settings info.Settings
	}

	dependencies struct {
		deps               deps.Deps
		depotBimap         bimap.Bimap[string, string]
		futureDependencies kv.Database[future.ID, string]
		taskDependencies   kv.Database[string, string]
	}

	parametrizationCodec codec.Codec[Parametrization, []byte]

	futures *future.Manager

	kernel struct {
		inputs  job2.Inputs
		outputs job2.Outputs
		db      job2.DB
	}

	settings Settings

	plugins []job2.Plugin
}

func Configure() Configurator {
	return Configurator{}
}

func (configurator Configurator) Info(storage record.Storage, settings info.Settings) Configurator {
	configurator.info.storage = storage
	configurator.info.settings = settings
	return configurator
}

func (configurator Configurator) Dependencies(
	deps deps.Deps,
	depotBimap bimap.Bimap[string, string],
	futureDependencies kv.Database[future.ID, string],
	taskDependencies kv.Database[string, string],
) Configurator {
	configurator.dependencies.deps = deps
	configurator.dependencies.depotBimap = depotBimap
	configurator.dependencies.futureDependencies = futureDependencies
	configurator.dependencies.taskDependencies = taskDependencies
	return configurator
}

func (configurator Configurator) Jobs(inputs job2.Inputs, outputs job2.Outputs, db job2.DB) Configurator {
	configurator.kernel.inputs = inputs
	configurator.kernel.outputs = outputs
	configurator.kernel.db = db
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

func (configurator Configurator) Settings(settings Settings) Configurator {
	configurator.settings = settings
	return configurator
}

func (configurator Configurator) Plugins(plugins ...job2.Plugin) Configurator {
	configurator.plugins = append(configurator.plugins, plugins...)
	return configurator
}

func (configurator Configurator) Compile() (*Kantoku, error) {
	dep := depot.New(
		configurator.dependencies.deps,
		configurator.dependencies.depotBimap,
		configurator.kernel.inputs,
	)

	kan := &Kantoku{
		info:    info.NewStorage(configurator.info.storage, configurator.info.settings),
		futures: configurator.futures,
		jobs:    job2.NewManager(dep, configurator.kernel.outputs, configurator.kernel.db),
		dependencies: DependencyManager{
			taskdep: taskdep.NewManager(configurator.dependencies.deps, configurator.dependencies.taskDependencies),
			futdep:  futdep.NewManager(configurator.dependencies.deps, configurator.dependencies.futureDependencies),
			depot:   dep,
		},
		parametrizationCodec: configurator.parametrizationCodec,
		settings:             configurator.settings,
	}

	plugins := []job2.Plugin{
		info.NewPlugin(kan.info),
		futdep.NewPlugin(kan.Dependencies().Futures(), kan.futures),
		taskdep.NewPlugin(kan.Dependencies().Tasks()),
		depot.NewPlugin(kan.Dependencies().Depot()),
	}
	plugins = append(plugins, configurator.plugins...)

	for _, plugin := range plugins {
		if err := kan.jobs.Use(plugin); err != nil {
			return nil, err
		}
	}

	return kan, nil
}

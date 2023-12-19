package alpha

import (
	"kantoku/common/data/identifier"
	"kantoku/common/data/record"
)

type Configurator struct {
	pool    Pool
	results ResultStorage
	storage Storage
	ids     identifier.Generator
	info    record.Storage
	runner  Runner
	plugins []Plugin
}

func (configurator Configurator) Pool(pool Pool) Configurator {
	configurator.pool = pool
	return configurator
}

func (configurator Configurator) Results(storage ResultStorage) Configurator {
	configurator.results = storage
	return configurator
}

func (configurator Configurator) Storage(storage Storage) Configurator {
	configurator.storage = storage
	return configurator
}

func (configurator Configurator) IdGenerator(generator identifier.Generator) Configurator {
	configurator.ids = generator
	return configurator
}

func (configurator Configurator) Runner(runner Runner) Configurator {
	configurator.runner = runner
	return configurator
}

func (configurator Configurator) Plugins(plugins []Plugin) Configurator {
	configurator.plugins = plugins
	return configurator
}

func (configurator Configurator) Compile() *Manager {
	return &Manager{
		pool:    configurator.pool,
		results: configurator.results,
		storage: configurator.storage,
		ids:     configurator.ids,
		info:    configurator.info,
		runner:  configurator.runner,
		plugins: configurator.plugins,
	}
}

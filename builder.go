package kantoku

import (
	"kantoku/common/codec"
	"kantoku/common/data/bimap"
	"kantoku/common/data/kv"
	"kantoku/common/data/record"
	"kantoku/framework/future"
	"kantoku/framework/plugins/depot"
	"kantoku/framework/plugins/depot/deps"
	"kantoku/framework/plugins/futdep"
	"kantoku/framework/plugins/info"
	"kantoku/framework/plugins/taskdep"
	"kantoku/kernel"
	"kantoku/kernel/platform"
)

type Builder struct {
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
	}

	settings Settings

	plugins []kernel.Plugin
}

func NewBuilder() Builder {
	return Builder{}
}

func (builder Builder) ConfigureInfo(storage record.Storage, settings info.Settings) Builder {
	builder.info.storage = storage
	builder.info.settings = settings
	return builder
}

func (builder Builder) ConfigureDeps(deps deps.Deps) Builder {
	builder.deps = deps
	return builder
}

func (builder Builder) ConfigureDepot(groupTaskBimap bimap.Bimap[string, string]) Builder {
	builder.depot.groupTaskBimap = groupTaskBimap
	return builder
}

func (builder Builder) ConfigurePlatform(platform platform.Platform[kernel.Task]) Builder {
	builder.platform = platform
	return builder
}

func (builder Builder) ConfigureParametrizationCodec(parametrizationCodec codec.Codec[Parametrization, []byte]) Builder {
	builder.parametrizationCodec = parametrizationCodec
	return builder
}

func (builder Builder) ConfigureFutures(futures *future.Manager) Builder {
	builder.futures = futures
	return builder
}

func (builder Builder) ConfigureFutdep(future2dependency kv.Database[future.ID, string]) Builder {
	builder.futdep.future2dependency = future2dependency
	return builder
}

func (builder Builder) ConfigureTaskdep(task2dependency kv.Database[string, string]) Builder {
	builder.taskdep.task2dependency = task2dependency
	return builder
}

func (builder Builder) ConfigureSettings(settings Settings) Builder {
	builder.settings = settings
	return builder
}

func (builder Builder) AddPlugins(plugins ...kernel.Plugin) Builder {
	builder.plugins = append(builder.plugins, plugins...)
	return builder
}

func (builder Builder) Build() *Kantoku {
	kan := &Kantoku{
		depot:                depot.New(builder.deps, builder.depot.groupTaskBimap, builder.platform.Inputs()),
		parametrizationCodec: builder.parametrizationCodec,
		info:                 info.NewStorage(builder.info.storage, builder.info.settings),
		futures:              builder.futures,
		futdep:               futdep.NewManager(builder.deps, builder.futdep.future2dependency),
		taskdep:              taskdep.NewManager(builder.deps, builder.taskdep.task2dependency),
		kernel:               kernel.New(builder.platform),
		settings:             builder.settings,
	}

	plugins := []kernel.Plugin{
		info.NewPlugin(kan.info),
		futdep.NewPlugin(kan.futdep),
		taskdep.NewPlugin(kan.taskdep),
	}
	plugins = append(plugins, builder.plugins...)
	for _, plugin := range plugins {
		kan.kernel.Register(plugin)
	}

	return kan
}

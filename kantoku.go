package kantoku

import (
	"context"
	"kantoku/common/codec"
	"kantoku/framework/future"
	"kantoku/framework/plugins/depot"
	"kantoku/framework/plugins/futdep"
	"kantoku/framework/plugins/info"
	"kantoku/framework/plugins/taskdep"
	"kantoku/kernel"
)

type Kantoku struct {
	depot                *depot.Depot
	parametrizationCodec codec.Codec[Parametrization, []byte]
	info                 *info.Storage
	futures              *future.Manager
	futdep               *futdep.Manager
	taskdep              *taskdep.Manager
	kernel               *kernel.Kernel
	settings             Settings
}

func (kantoku *Kantoku) Spawn(ctx context.Context, spec Spec) (kernel.Result, error) {
	payload, err := kantoku.parametrizationCodec.Encode(spec.parametrization)
	if err != nil {
		return kernel.Result{}, err
	}
	var options []kernel.Option

	options = append(options, withParametrization(spec.parametrization))
	options = append(options, spec.opts...)

	if kantoku.settings.AutoInputDependencies {
		options = append(options, AutoInputDeps())
	}

	return kantoku.kernel.Spawn(ctx, kernel.Describe(spec.typ, payload).With(options...))
}

func (kantoku *Kantoku) Task(id string) Task {
	return Task{
		id:      id,
		kantoku: kantoku,
	}
}

func (kantoku *Kantoku) Info() *info.Storage {
	return kantoku.info
}

func (kantoku *Kantoku) Depot() *depot.Depot {
	return kantoku.depot
}

func (kantoku *Kantoku) Futures() *future.Manager {
	return kantoku.futures
}

func (kantoku *Kantoku) Kernel() *kernel.Kernel {
	return kantoku.kernel
}

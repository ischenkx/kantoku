package deprecated

import (
	"github.com/samber/lo"
	"kantoku/common/data/future"
	"kantoku/framework/job"
	"kantoku/framework/plugins/depot"
	"kantoku/framework/plugins/info"
)

type Spec struct {
	typ             string
	parametrization Parametrization
	opts            []job.Option
}

func Describe(typ string) Spec {
	return Spec{typ: typ}
}

func (spec Spec) WithContext(ctx string) Spec {
	return spec.WithMeta("context", ctx)
}

func (spec Spec) WithDeps(dependencies ...string) Spec {
	opts := lo.Map(dependencies, func(dep string, _ int) job.Option { return depot.Dependency(dep) })
	return spec.With(opts...)
}

func (spec Spec) WithMeta(key string, value any) Spec {
	return spec.With(info.WithEntry(key, value))
}

func (spec Spec) WithStatic(static []byte) Spec {
	spec.parametrization.Static = static
	return spec
}

func (spec Spec) WithInputs(ids ...future.ID) Spec {
	spec.parametrization.Inputs = ids
	return spec
}

func (spec Spec) WithOutputs(ids ...future.ID) Spec {
	spec.parametrization.Outputs = ids
	return spec
}

func (spec Spec) With(opts ...job.Option) Spec {
	spec.opts = append(spec.opts, opts...)
	return spec
}

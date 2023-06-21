package kantoku

import (
	"github.com/samber/lo"
	"kantoku/framework/future"
	"kantoku/framework/plugins/depot"
	"kantoku/framework/plugins/meta"
	"kantoku/kernel"
)

type Spec struct {
	typ             string
	parametrization Parametrization
	opts            []kernel.Option
}

func Describe(typ string) Spec {
	return Spec{typ: typ}
}

func (spec Spec) WithContext(ctx string) Spec {
	return spec.WithMeta("context", ctx)
}

func (spec Spec) WithDeps(dependencies ...string) Spec {
	opts := lo.Map(dependencies, func(dep string, _ int) kernel.Option { return depot.Dependency(dep) })
	return spec.WithOptions(opts...)
}

func (spec Spec) WithMeta(key string, value any) Spec {
	return spec.WithOptions(meta.WithEntry(key, value))
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

func (spec Spec) WithOptions(opts ...kernel.Option) Spec {
	spec.opts = append(spec.opts, opts...)
	return spec
}

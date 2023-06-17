package kantoku

import (
	"kantoku/framework/plugins/futdep"
	"kantoku/kernel"
)

func withParametrization(param Parametrization) kernel.Option {
	return func(ctx *kernel.Context) error {
		ctx.Data().Set("param", param)
		return nil
	}
}

func autoInputDeps(ctx *kernel.Context) error {
	param := GetParametrization(ctx)
	for _, id := range param.Inputs {
		if err := futdep.Dep(id)(ctx); err != nil {
			return err
		}
	}
	return nil
}

func AutoInputDeps() kernel.Option {
	return autoInputDeps
}

func GetParametrization(ctx *kernel.Context) Parametrization {
	return ctx.Data().GetWithDefault("param", Parametrization{}).(Parametrization)
}

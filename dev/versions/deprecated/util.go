package deprecated

import (
	job2 "kantoku/framework/job"
	"kantoku/framework/plugins/futdep"
)

func withParametrization(param Parametrization) job2.Option {
	return func(ctx *job2.Context) error {
		ctx.Data().Set("param", param)
		return nil
	}
}

func autoInputDeps(ctx *job2.Context) error {
	param := GetParametrization(ctx)
	for _, id := range param.Inputs {
		if err := futdep.Dep(id)(ctx); err != nil {
			return err
		}
	}
	return nil
}

func AutoInputDeps() job2.Option {
	return autoInputDeps
}

func GetParametrization(ctx *job2.Context) Parametrization {
	return ctx.Data().GetWithDefault("param", Parametrization{}).(Parametrization)
}

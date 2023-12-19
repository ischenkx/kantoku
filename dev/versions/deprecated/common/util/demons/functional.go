package demons

import (
	"context"
	"kantoku/framework/infra"
)

func Functional(name string, fn func(ctx context.Context) error) infra.Demon {
	return infra.Demon{
		Type:      "FUNCTIONAL",
		Name:      name,
		Parameter: fn,
	}
}

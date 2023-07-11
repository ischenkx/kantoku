package demons

import (
	"context"
	"kantoku/framework/infra/demon"
)

func Functional(name string, fn func(ctx context.Context) error) demon.Demon {
	return demon.Demon{
		Type:      "FUNCTIONAL",
		Name:      name,
		Parameter: fn,
	}
}

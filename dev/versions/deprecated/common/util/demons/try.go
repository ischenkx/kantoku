package demons

import (
	"kantoku/framework/infra"
)

func TryProvider(obj any) infra.Provider {
	if provider, ok := obj.(infra.Provider); ok {
		return provider
	}
	return EmptyProvider
}

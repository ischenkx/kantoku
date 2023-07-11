package demons

import "kantoku/framework/infra/demon"

func TryProvider(obj any) demon.Provider {
	if provider, ok := obj.(demon.Provider); ok {
		return provider
	}
	return EmptyProvider
}

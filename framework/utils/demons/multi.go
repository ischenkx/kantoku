package demons

import (
	"context"
	"kantoku/framework/infra/demon"
)

type Multi []demon.Provider

var EmptyProvider = Multi{}

func (m Multi) Demons(ctx context.Context) (result []demon.Demon) {
	for _, provider := range m {
		result = append(result, provider.Demons(ctx)...)
	}

	return result
}

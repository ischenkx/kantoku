package demons

import (
	infra2 "kantoku/framework/infra"
)

type Multi []infra2.Provider

var EmptyProvider = Multi{}

func (m Multi) Demons() (result []infra2.Demon) {
	for _, provider := range m {
		result = append(result, provider.Demons()...)
	}

	return result
}

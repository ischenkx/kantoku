package argument

import (
	"context"
	"errors"
	"fmt"
)

type Codecs map[string]Codec

type Manager struct {
	codecs         Codecs
	defaultEncoder Codec
}

func NewManager(codecs Codecs, defaultEncoder Codec) *Manager {
	return &Manager{codecs: codecs, defaultEncoder: defaultEncoder}
}

func (manager *Manager) Encode(ctx context.Context, raw any) (Argument, error) {
	if selfEncoder, ok := raw.(SelfEncoder); ok {
		return selfEncoder.Encode(ctx)
	}

	if manager.defaultEncoder == nil {
		return Argument{}, errors.New("no default encoder provided")
	}

	return manager.defaultEncoder.Encode(ctx, raw)
}

func (manager *Manager) Decode(ctx context.Context, argument Argument) (any, error) {
	codec, ok := manager.codecs[argument.Type]
	if !ok {
		return nil, fmt.Errorf("no codec for a given type '%s'", argument.Type)
	}

	return codec.Decode(ctx, argument)
}

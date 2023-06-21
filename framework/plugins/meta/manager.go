package meta

import (
	"context"
	"kantoku/common/codec"
	"kantoku/common/data/kv"
)

type Storage kv.Getter[string, RawMeta]

type Manager struct {
	storage Storage
	codec   codec.Dynamic[[]byte]
}

func NewManager(storage Storage, codec codec.Dynamic[[]byte]) *Manager {
	return &Manager{storage: storage, codec: codec}
}

func (manager *Manager) Get(ctx context.Context, id string) (Meta, error) {
	raw, err := manager.storage.Get(ctx, id)
	if err != nil {
		return Meta{}, err
	}

	return Meta{
		raw:   raw,
		codec: manager.codec,
	}, nil
}

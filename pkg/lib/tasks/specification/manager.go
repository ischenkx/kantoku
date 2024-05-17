package specification

import (
	"context"
	"encoding/json"
	"fmt"
)

type BinaryStorage interface {
	Get(ctx context.Context, id string) ([]byte, error)
	GetAll(ctx context.Context) ([][]byte, error)
	Add(ctx context.Context, id string, data []byte) error
	Remove(ctx context.Context, id string) error
}

type JsonStorage[T any] struct {
	Raw    BinaryStorage
	IdFunc func(T) string
}

func zeroValue[T any]() T {
	var zero T
	return zero
}

func (storage *JsonStorage[T]) Get(ctx context.Context, id string) (T, error) {
	rawSpec, err := storage.Raw.Get(ctx, id)
	if err != nil {
		return zeroValue[T](), err
	}

	var data T
	if err := json.Unmarshal(rawSpec, &data); err != nil {
		return zeroValue[T](), fmt.Errorf("failed to unmarshal: %w", err)
	}

	return data, nil
}

func (storage *JsonStorage[T]) GetAll(ctx context.Context) ([]T, error) {
	rawItems, err := storage.Raw.GetAll(ctx)
	if err != nil {
		return nil, err
	}

	items := make([]T, 0, len(rawItems))
	for _, rawItem := range rawItems {
		var item T
		if err := json.Unmarshal(rawItem, &item); err != nil {
			return nil, fmt.Errorf("failed to unmarshal specification: %w", err)
		}

		items = append(items, item)
	}

	return items, nil
}

func (storage *JsonStorage[T]) Add(ctx context.Context, item T) error {
	encodedItem, err := json.Marshal(item)
	if err != nil {
		return fmt.Errorf("failed to marshal: %w", err)
	}

	return storage.Raw.Add(ctx, storage.IdFunc(item), encodedItem)
}

func (storage *JsonStorage[T]) Remove(ctx context.Context, id string) error {
	return storage.Raw.Remove(ctx, id)
}

type Manager struct {
	specifications *JsonStorage[Specification]
	types          *JsonStorage[TypeWithID]
}

func NewManager(specificationBinaryStorage, typeBinaryStorage BinaryStorage) *Manager {
	return &Manager{
		specifications: &JsonStorage[Specification]{
			Raw: specificationBinaryStorage,
			IdFunc: func(specification Specification) string {
				return specification.ID
			},
		},
		types: &JsonStorage[TypeWithID]{
			Raw: typeBinaryStorage,
			IdFunc: func(t TypeWithID) string {
				return t.ID
			},
		},
	}
}

func (manager *Manager) Specifications() *JsonStorage[Specification] {
	return manager.specifications
}

func (manager *Manager) Types() *JsonStorage[TypeWithID] {
	return manager.types
}

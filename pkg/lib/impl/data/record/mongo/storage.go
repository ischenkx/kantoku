package mongorec

import (
	"context"
	"fmt"
	"github.com/ischenkx/kantoku/pkg/common/data/codec"
	"github.com/ischenkx/kantoku/pkg/common/data/record"
	"go.mongodb.org/mongo-driver/mongo"
)

var _ record.Storage[int] = (*Storage[int])(nil)

type Storage[Item any] struct {
	collection *mongo.Collection
	codec      codec.Codec[Item, record.R]
}

func New[Item any](collection *mongo.Collection, codec codec.Codec[Item, record.R]) *Storage[Item] {
	return &Storage[Item]{collection: collection, codec: codec}
}

func (storage *Storage[Item]) Insert(ctx context.Context, item Item) error {
	rec, err := storage.codec.Encode(item)
	if err != nil {
		return fmt.Errorf("failed to encode: %w", err)
	}
	if len(rec) == 0 {
		return nil
	}
	_, err = storage.collection.InsertOne(ctx, record2bson(rec))
	if err != nil {
		return fmt.Errorf("failed to insert: %s", err)
	}

	return nil
}

func (storage *Storage[Item]) Filter(record record.Record) record.Set[Item] {
	return newSet(storage).Filter(record)
}

func (storage *Storage[Item]) Erase(ctx context.Context) error {
	return newSet(storage).Erase(ctx)
}

func (storage *Storage[Item]) Update(ctx context.Context, update, upsert record.R) error {
	return newSet(storage).Update(ctx, update, upsert)
}

func (storage *Storage[Item]) Distinct(keys ...string) record.Cursor[Item] {
	return newSet(storage).Distinct(keys...)
}

func (storage *Storage[Item]) Cursor() record.Cursor[Item] {
	return newSet(storage).Cursor()
}

package mongorec

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/mongo"
	"kantoku/common/data/record"
)

var _ record.Storage = (*Storage)(nil)

type Storage struct {
	collection *mongo.Collection
}

func New(collection *mongo.Collection) *Storage {
	return &Storage{collection: collection}
}

func (storage *Storage) Insert(ctx context.Context, rec record.Record) error {
	if len(rec) == 0 {
		return nil
	}
	_, err := storage.collection.InsertOne(ctx, record2bson(rec))
	if err != nil {
		return fmt.Errorf("failed to insert: %s", err)
	}

	return nil
}

func (storage *Storage) Filter(entries ...record.Entry) record.Set {
	return newSet(storage).Filter(entries...)
}

func (storage *Storage) Erase(ctx context.Context) error {
	return newSet(storage).Erase(ctx)
}

func (storage *Storage) Update(ctx context.Context, update, upsert record.R) error {
	return newSet(storage).Update(ctx, update, upsert)
}

func (storage *Storage) Distinct(keys ...string) record.Cursor[record.Record] {
	return newSet(storage).Distinct(keys...)
}

func (storage *Storage) Cursor() record.Cursor[record.Record] {
	return newSet(storage).Cursor()
}

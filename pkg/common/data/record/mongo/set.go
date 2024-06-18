package mongorec

import (
	"context"
	"fmt"
	"github.com/ischenkx/kantoku/pkg/common/data/record"
	"github.com/ischenkx/kantoku/pkg/common/data/record/ops"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var _ record.Set[int] = (*Set[int])(nil)

type Set[Item any] struct {
	storage *Storage[Item]
	filter  record.R
}

func newSet[Item any](storage *Storage[Item]) Set[Item] {
	return Set[Item]{storage: storage}
}

func (set Set[Item]) Filter(rec record.R) record.Set[Item] {
	newFilter := record.R{}
	for key, value := range set.filter {
		newFilter[key] = value
	}

	for key, value := range rec {
		if oldKey, ok := newFilter[key]; ok {
			value = ops.And(value, oldKey)
		}
		newFilter[key] = value
	}

	set.filter = newFilter
	return set
}

func (set Set[Item]) Distinct(keys ...string) record.Cursor[Item] {
	return DistinctCursor[Item]{
		skip:    0,
		limit:   0,
		filter:  set.filter,
		keys:    keys,
		storage: set.storage,
	}
}

func (set Set[Item]) Erase(ctx context.Context) error {
	filter, err := makeRecordFilter(set.filter)
	if err != nil {
		return fmt.Errorf("failed to make a filter: %w", err)
	}
	_, err = set.storage.collection.DeleteMany(ctx, filter)
	if err != nil {
		return fmt.Errorf("failed to delete: %s", err)
	}

	return nil
}

func (set Set[Item]) Update(ctx context.Context, update, upsert record.R) error {
	bsonUpdate := bson.M{"$set": record2bson(unwrapUpdateRecord(update))}

	if upsert != nil {
		bsonUpsert := record2bson(upsert)
		bsonSetter := bsonUpdate["$set"].(bson.M)
		for key := range bsonSetter {
			delete(bsonUpsert, key)
		}
		bsonUpdate["$setOnInsert"] = bsonUpsert
	}

	filter, err := makeRecordFilter(set.filter)
	if err != nil {
		return fmt.Errorf("failed to make a filter: %w", err)
	}

	_, err = set.storage.collection.UpdateMany(ctx,
		filter,
		bsonUpdate,
		options.
			Update().
			SetUpsert(upsert != nil),
	)

	if err != nil {
		return fmt.Errorf("failed to update many: %s", err)
	}

	return nil
}

func (set Set[Item]) Cursor() record.Cursor[Item] {
	return FilterCursor[Item]{
		skip:    0,
		limit:   0,
		filter:  set.filter,
		storage: set.storage,
	}
}

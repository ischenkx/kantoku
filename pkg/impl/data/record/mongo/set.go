package mongorec

import (
	"context"
	"fmt"
	"github.com/ischenkx/kantoku/pkg/common/data/record"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var _ record.Set = (*Set)(nil)

type Set struct {
	storage *Storage
	filter  record.R
}

func newSet(storage *Storage) Set {
	return Set{storage: storage}
}

func (set Set) Filter(rec record.R) record.Set {
	set.filter = rec
	return set
}

func (set Set) Distinct(keys ...string) record.Cursor[record.Record] {
	return DistinctCursor{
		skip:    0,
		limit:   0,
		filter:  set.filter,
		keys:    keys,
		storage: set.storage,
	}
}

func (set Set) Erase(ctx context.Context) error {
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

func (set Set) Update(ctx context.Context, update, upsert record.R) error {
	bsonUpdate := bson.M{"$set": record2bson(update)}

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

func (set Set) Cursor() record.Cursor[record.Record] {
	return FilterCursor{
		skip:    0,
		limit:   0,
		filter:  set.filter,
		storage: set.storage,
	}
}

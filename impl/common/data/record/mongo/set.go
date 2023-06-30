package mongorec

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
	"kantoku/common/data/record"
)

var _ record.Set = (*Set)(nil)

type Set struct {
	storage *Storage
	filters [][]record.Entry
}

func newSet(storage *Storage) Set {
	return Set{storage: storage}
}

func (set Set) Filter(entries ...record.Entry) record.Set {
	set.filters = append(set.filters, entries)
	return set
}

func (set Set) Distinct(keys ...string) record.Cursor[record.Record] {
	return DistinctCursor{
		skip:    0,
		limit:   0,
		filters: set.filters,
		keys:    keys,
		storage: set.storage,
	}
}

func (set Set) Erase(ctx context.Context) error {
	_, err := set.storage.collection.DeleteMany(ctx, makeFilter(set.filters))
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

	_, err := set.storage.collection.UpdateMany(ctx, makeFilter(set.filters), bsonUpdate, options.Update().SetUpsert(upsert != nil))

	if err != nil {
		return fmt.Errorf("failed to update many: %s", err)
	}

	return nil
}

func (set Set) Cursor() record.Cursor[record.Record] {
	return FilterCursor{
		skip:    0,
		limit:   0,
		filters: set.filters,
		storage: set.storage,
	}
}

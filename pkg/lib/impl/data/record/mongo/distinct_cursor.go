package mongorec

import (
	"context"
	"fmt"
	"github.com/ischenkx/kantoku/pkg/common/data/record"
	"github.com/samber/lo"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

var _ record.Cursor[int] = DistinctCursor[int]{}

type DistinctCursor[Item any] struct {
	skip    int
	limit   int
	filter  record.R
	sorters []record.Sorter
	masks   []record.Mask
	keys    []string
	storage *Storage[Item]
}

func (cursor DistinctCursor[Item]) Skip(num int) record.Cursor[Item] {
	cursor.skip += num
	return cursor
}

func (cursor DistinctCursor[Item]) Limit(num int) record.Cursor[Item] {
	cursor.limit = num
	return cursor
}

func (cursor DistinctCursor[Item]) Mask(masks ...record.Mask) record.Cursor[Item] {
	cursor.masks = append(cursor.masks, masks...)
	return cursor
}

func (cursor DistinctCursor[Item]) Sort(sorters ...record.Sorter) record.Cursor[Item] {
	cursor.sorters = sorters
	return cursor
}

func (cursor DistinctCursor[Item]) Iter() record.Iter[Item] {
	return &DistinctIter[Item]{DistinctCursor: cursor}
}

func (cursor DistinctCursor[Item]) Count(ctx context.Context) (int, error) {
	if len(cursor.keys) == 0 {
		return 0, nil
	}

	filter, err := makeRecordFilter(cursor.filter)
	if err != nil {
		return 0, fmt.Errorf("failed to make a filter: %w", err)
	}

	pipeline := bson.A{
		bson.M{"$match": filter},
		bson.M{
			"$group": bson.M{
				"_id": lo.SliceToMap[string, string, any](
					cursor.keys,
					func(key string) (string, any) {
						return key, fmt.Sprintf("$%s", key)
					},
				),
			},
		},
		bson.M{"$group": bson.M{"_id": nil, "totalCount": bson.M{"$sum": 1}}},
	}
	result, err := cursor.storage.collection.Aggregate(ctx, pipeline)
	if err != nil {
		return 0, fmt.Errorf("failed to aggregate: %s", err)
	}
	defer result.Close(ctx)

	var doc struct {
		TotalCount int `bson:"totalCount"`
	}

	if !result.Next(ctx) {
		return 0, nil
	}

	if err := result.Decode(&doc); err != nil {
		return 0, fmt.Errorf("failed to decode the document: %s", err)
	}

	return doc.TotalCount, nil
}

type DistinctIter[Item any] struct {
	DistinctCursor[Item]
	mongoCursor *mongo.Cursor
}

func (iter *DistinctIter[Item]) Close(ctx context.Context) error {
	if iter.mongoCursor == nil {
		return nil
	}
	return iter.mongoCursor.Close(ctx)
}

func (iter *DistinctIter[Item]) Next(ctx context.Context) (Item, error) {
	var zero Item

	if len(iter.keys) == 0 {
		return zero, record.ErrIterEmpty
	}

	cursor, err := iter.getCursor(ctx)
	if err != nil {
		return zero, fmt.Errorf("failed to make a mongo cursor: %s", err)
	}

	if !cursor.Next(ctx) {
		return zero, record.ErrIterEmpty
	}

	var doc struct {
		ID bson.M `bson:"_id"`
	}

	err = cursor.Decode(&doc)
	if err != nil {
		return zero, fmt.Errorf("failed to decode the received data: %s", err)
	}

	rec := bson2record(doc.ID)
	for _, mask := range iter.masks {
		if mask.Operation == record.IncludeMask {
			if _, ok := rec[mask.PropertyPattern]; !ok {
				rec[mask.PropertyPattern] = nil
			}
		}
	}

	item, err := iter.DistinctCursor.storage.codec.Decode(rec)
	if err != nil {
		return zero, fmt.Errorf("failed to decode: %w", err)
	}

	return item, nil
}

func (iter *DistinctIter[Item]) getCursor(ctx context.Context) (*mongo.Cursor, error) {
	if iter.mongoCursor != nil {
		return iter.mongoCursor, nil
	}

	filter, err := makeRecordFilter(iter.filter)
	if err != nil {
		return nil, fmt.Errorf("failed to make a filter: %w", err)
	}

	pipeline := bson.A{
		bson.M{"$match": filter},
		bson.M{
			"$group": bson.M{
				"_id": lo.SliceToMap[string, string, any](
					iter.keys,
					func(key string) (string, any) {
						return key, bson.M{"$ifNull": bson.A{fmt.Sprintf("$%s", key), nil}}
					},
				),
			},
		},
	}

	if len(iter.sorters) > 0 {
		sort := bson.D{}
		for _, sorter := range iter.sorters {
			switch sorter.Ordering {
			case record.Asc:
				sort = append(sort, bson.E{fmt.Sprintf("_id.%s", sorter.Key), 1})
			case record.Desc:
				sort = append(sort, bson.E{fmt.Sprintf("_id.%s", sorter.Key), -1})
			}
		}

		pipeline = append(pipeline, bson.M{"$sort": sort})
	}

	if len(iter.masks) > 0 {
		projection := bson.M{}
		for _, mask := range iter.masks {
			switch mask.Operation {
			case record.IncludeMask:
				projection[fmt.Sprintf("_id.%s", mask.PropertyPattern)] = 1
			case record.ExcludeMask:
				projection[fmt.Sprintf("_id.%s", mask.PropertyPattern)] = 0
			}
		}
		pipeline = append(pipeline, bson.M{"$project": projection})
	}

	if iter.skip > 0 {
		pipeline = append(pipeline, bson.M{"$skip": iter.skip})
	}

	if iter.limit > 0 {
		pipeline = append(pipeline, bson.M{"$limit": iter.limit})
	}

	cursor, err := iter.storage.collection.Aggregate(ctx, pipeline)

	if err != nil {
		return nil, fmt.Errorf("failed to aggregate: %s", err)
	}

	iter.mongoCursor = cursor
	return iter.mongoCursor, nil
}

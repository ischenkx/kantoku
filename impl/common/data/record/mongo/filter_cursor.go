package mongorec

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"kantoku/common/data/record"
)

type FilterCursor struct {
	skip    int
	limit   int
	filters [][]record.Entry
	sorters []record.Sorter
	storage *Storage
	masks   []record.Mask
}

func (f FilterCursor) Skip(num int) record.Cursor[record.Record] {
	f.skip += num
	return f
}

func (f FilterCursor) Limit(num int) record.Cursor[record.Record] {
	f.limit = num
	return f
}

func (f FilterCursor) Mask(masks ...record.Mask) record.Cursor[record.Record] {
	f.masks = append(f.masks, masks...)
	return f
}

func (f FilterCursor) Sort(sorters ...record.Sorter) record.Cursor[record.Record] {
	f.sorters = sorters
	return f
}

func (f FilterCursor) Iter() record.Iter[record.Record] {
	return &FilterIter{FilterCursor: f}
}

func (f FilterCursor) Count(ctx context.Context) (int, error) {
	num, err := f.storage.collection.CountDocuments(ctx, makeFilter(f.filters))
	if err != nil {
		return 0, fmt.Errorf("failed to count documents: %s", err)
	}

	return int(num), nil
}

type FilterIter struct {
	FilterCursor
	mongoCursor *mongo.Cursor
}

func (iter *FilterIter) Close(ctx context.Context) error {
	if iter.mongoCursor == nil {
		return nil
	}
	return iter.mongoCursor.Close(ctx)
}

func (iter *FilterIter) Next(ctx context.Context) (record.Record, error) {
	cursor, err := iter.getCursor(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to make a mongo cursor: %s", err)
	}

	if !cursor.Next(ctx) {
		return nil, record.ErrIterEmpty
	}

	var doc bson.M
	err = cursor.Decode(&doc)
	if err != nil {
		return nil, fmt.Errorf("failed to decode the received data: %s", err)
	}
	delete(doc, "_id")

	rec := bson2record(doc)
	for _, mask := range iter.masks {
		if mask.Operation == record.IncludeMask {
			if _, ok := rec[mask.PropertyPattern]; !ok {
				rec[mask.PropertyPattern] = nil
			}
		}
	}

	return rec, nil
}

func (iter *FilterIter) getCursor(ctx context.Context) (*mongo.Cursor, error) {
	if iter.mongoCursor != nil {
		return iter.mongoCursor, nil
	}

	opts := options.Find()
	if len(iter.sorters) > 0 {
		sort := bson.D{}
		for _, sorter := range iter.sorters {
			switch sorter.Ordering {
			case record.ASC:
				sort = append(sort, bson.E{sorter.Key, 1})
			case record.DESC:
				sort = append(sort, bson.E{sorter.Key, -1})
			}
		}
		opts.SetSort(sort)
	}

	if len(iter.masks) > 0 {
		projection := bson.M{}
		for _, mask := range iter.masks {
			switch mask.Operation {
			case record.IncludeMask:
				projection[mask.PropertyPattern] = 1
			case record.ExcludeMask:
				projection[mask.PropertyPattern] = 0
			}
		}
		opts.SetProjection(projection)
	}

	if iter.skip > 0 {
		opts.SetSkip(int64(iter.skip))
	}

	if iter.limit > 0 {
		opts.SetLimit(int64(iter.limit))
	}

	cursor, err := iter.storage.collection.Find(ctx, makeFilter(iter.filters), opts)
	if err != nil {
		return nil, fmt.Errorf("failed to find: %s", err)
	}

	iter.mongoCursor = cursor
	return iter.mongoCursor, nil
}

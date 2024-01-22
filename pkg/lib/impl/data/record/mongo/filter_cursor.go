package mongorec

import (
	"context"
	"fmt"
	"github.com/ischenkx/kantoku/pkg/common/data/record"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var _ record.Cursor[int] = FilterCursor[int]{}

type FilterCursor[Item any] struct {
	skip    int
	limit   int
	filter  record.Record
	sorters []record.Sorter
	storage *Storage[Item]
	masks   []record.Mask
}

func (f FilterCursor[Item]) Skip(num int) record.Cursor[Item] {
	f.skip += num
	return f
}

func (f FilterCursor[Item]) Limit(num int) record.Cursor[Item] {
	f.limit = num
	return f
}

func (f FilterCursor[Item]) Mask(masks ...record.Mask) record.Cursor[Item] {
	f.masks = append(f.masks, masks...)
	return f
}

func (f FilterCursor[Item]) Sort(sorters ...record.Sorter) record.Cursor[Item] {
	f.sorters = sorters
	return f
}

func (f FilterCursor[Item]) Iter() record.Iter[Item] {
	return &FilterIter[Item]{FilterCursor: f}
}

func (f FilterCursor[Item]) Count(ctx context.Context) (int, error) {
	filter, err := makeRecordFilter(f.filter)
	if err != nil {
		return 0, fmt.Errorf("failed to make a filter: %w", err)
	}
	num, err := f.storage.collection.CountDocuments(ctx, filter)
	if err != nil {
		return 0, fmt.Errorf("failed to count documents: %s", err)
	}

	return int(num), nil
}

type FilterIter[Item any] struct {
	FilterCursor[Item]
	mongoCursor *mongo.Cursor
}

func (iter *FilterIter[Item]) Close(ctx context.Context) error {
	if iter.mongoCursor == nil {
		return nil
	}
	return iter.mongoCursor.Close(ctx)
}

func (iter *FilterIter[Item]) Next(ctx context.Context) (Item, error) {
	var zero Item
	cursor, err := iter.getCursor(ctx)
	if err != nil {
		return zero, fmt.Errorf("failed to make a mongo cursor: %s", err)
	}

	if !cursor.Next(ctx) {
		return zero, record.ErrIterEmpty
	}

	var doc bson.M
	err = cursor.Decode(&doc)
	if err != nil {
		return zero, fmt.Errorf("failed to decode the received data: %s", err)
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

	item, err := iter.FilterCursor.storage.codec.Decode(rec)
	if err != nil {
		return zero, fmt.Errorf("failed to decode: %w", err)
	}

	return item, nil
}

func (iter *FilterIter[Item]) getCursor(ctx context.Context) (*mongo.Cursor, error) {
	if iter.mongoCursor != nil {
		return iter.mongoCursor, nil
	}

	opts := options.Find()
	if len(iter.sorters) > 0 {
		sort := bson.D{}
		for _, sorter := range iter.sorters {
			switch sorter.Ordering {
			case record.Asc:
				sort = append(sort, bson.E{sorter.Key, 1})
			case record.Desc:
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

	filter, err := makeRecordFilter(iter.filter)
	if err != nil {
		return nil, fmt.Errorf("failed to make a filter: %w", err)
	}

	cursor, err := iter.storage.collection.Find(ctx, filter, opts)
	if err != nil {
		return nil, fmt.Errorf("failed to find: %s", err)
	}

	iter.mongoCursor = cursor
	return iter.mongoCursor, nil
}

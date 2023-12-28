package client

import (
	"context"
	"errors"
	"fmt"
	"github.com/ischenkx/kantoku/pkg/common/data/record"
	"github.com/ischenkx/kantoku/pkg/extensions/web/oas"
	"github.com/samber/lo"
	"net/http"
)

type recordStorage struct {
	httpClient oas.ClientWithResponsesInterface
}

func (storage recordStorage) Insert(ctx context.Context, record record.Record) error {
	res, err := storage.httpClient.PostTasksInfoInsertWithResponse(ctx, record)
	if err != nil {
		return fmt.Errorf("failed to make an http request: %w", err)
	}

	code := res.StatusCode()

	switch code {
	case http.StatusOK:
		return nil
	case http.StatusInternalServerError:
		return fmt.Errorf("server failure: %s", res.JSON500.Message)
	default:
		return fmt.Errorf("unexpected response code: %d", code)
	}
}

func (storage recordStorage) Filter(rec record.Record) record.Set {
	return recordSet{
		httpClient: storage.httpClient,
		filter:     rec,
	}
}

func (storage recordStorage) Erase(ctx context.Context) error {
	return recordSet{httpClient: storage.httpClient}.Erase(ctx)
}

func (storage recordStorage) Update(ctx context.Context, update, upsert record.R) error {
	return recordSet{httpClient: storage.httpClient}.Update(ctx, update, upsert)
}

func (storage recordStorage) Distinct(keys ...string) record.Cursor[record.Record] {
	return recordSet{httpClient: storage.httpClient}.Distinct(keys...)
}

func (storage recordStorage) Cursor() record.Cursor[record.Record] {
	return recordSet{httpClient: storage.httpClient}.Cursor()
}

type recordSet struct {
	httpClient oas.ClientWithResponsesInterface

	filter record.Record
}

func (set recordSet) Filter(rec record.Record) record.Set {
	set.filter = rec
	return set
}

func (set recordSet) Erase(ctx context.Context) error {
	var body oas.PostTasksInfoEraseJSONRequestBody

	if set.filter != nil {
		body.Filter = (*oas.InfoFilter)(&set.filter)
	}

	res, err := set.httpClient.PostTasksInfoEraseWithResponse(ctx, body)
	if err != nil {
		return fmt.Errorf("failed to make an http request: %w", err)
	}

	code := res.StatusCode()

	switch code {
	case http.StatusOK:
		return nil
	case http.StatusInternalServerError:
		return fmt.Errorf("server failure: %s", res.JSON500.Message)
	default:
		return fmt.Errorf("unexpected response code: %d", code)
	}
}

func (set recordSet) Update(ctx context.Context, update, upsert record.R) error {
	body := oas.PostTasksInfoUpdateJSONRequestBody{
		Filter: set.filter,
		Update: update,
		Upsert: nil,
	}
	if upsert != nil {
		body.Upsert = (*oas.TaskInfo)(&upsert)
	}

	res, err := set.httpClient.PostTasksInfoUpdateWithResponse(ctx, body)
	if err != nil {
		return fmt.Errorf("failed to make an http request: %w", err)
	}

	code := res.StatusCode()

	switch code {
	case http.StatusOK:
		return nil
	case http.StatusInternalServerError:
		return fmt.Errorf("server failure: %s", res.JSON500.Message)
	default:
		return fmt.Errorf("unexpected response code: %d", code)
	}
}

func (set recordSet) Distinct(keys ...string) record.Cursor[record.Record] {
	cursor := set.makeCursor()
	cursor.distinctKeys = keys

	return cursor
}

func (set recordSet) Cursor() record.Cursor[record.Record] {
	return set.makeCursor()
}

func (set recordSet) makeCursor() recordCursor {
	return recordCursor{
		filter:     set.filter,
		httpClient: set.httpClient,
		limit:      -1,
	}
}

type recordCursor struct {
	httpClient oas.ClientWithResponsesInterface

	filter       record.Record
	distinctKeys []string
	masks        []record.Mask
	sorters      []record.Sorter
	skip         int
	limit        int
}

func (cursor recordCursor) Skip(n int) record.Cursor[record.Record] {
	cursor.skip = n
	return cursor
}

func (cursor recordCursor) Limit(n int) record.Cursor[record.Record] {
	cursor.limit = n
	return cursor
}

func (cursor recordCursor) Mask(masks ...record.Mask) record.Cursor[record.Record] {
	cursor.masks = append(cursor.masks, masks...)
	return cursor
}

func (cursor recordCursor) Sort(sorters ...record.Sorter) record.Cursor[record.Record] {
	cursor.sorters = sorters
	return cursor
}

func (cursor recordCursor) Iter() record.Iter[record.Record] {
	return &iter{
		cursor: cursor,
		index:  0,
		data:   nil,
		loaded: false,
		closed: false,
	}
}

func (cursor recordCursor) Count(ctx context.Context) (int, error) {
	body := oas.PostTasksInfoCountJSONRequestBody{
		Cursor: cursor.makeCursor(),
		Filter: cursor.makeFilter(),
	}

	res, err := cursor.httpClient.PostTasksInfoCountWithResponse(ctx, body)
	if err != nil {
		return 0, fmt.Errorf("failed to make an http request: %w", err)
	}

	code := res.StatusCode()

	switch code {
	case http.StatusOK:
		return *res.JSON200, nil
	case http.StatusInternalServerError:
		return 0, fmt.Errorf("server failure: %s", res.JSON500.Message)
	default:
		return 0, fmt.Errorf("unexpected response code: %d", code)
	}
}

func (cursor recordCursor) makeFilter() *oas.InfoFilter {
	if cursor.filter == nil {
		return nil
	}

	return (*oas.InfoFilter)(&cursor.filter)
}

func (cursor recordCursor) makeCursor() *oas.InfoCursor {
	infoCursor := &oas.InfoCursor{}

	if cursor.skip > 0 {
		infoCursor.Skip = &cursor.skip
	}

	if cursor.limit >= 0 {
		infoCursor.Skip = &cursor.limit
	}

	if len(cursor.distinctKeys) > 0 {
		infoCursor.Distinct = &cursor.distinctKeys
	}

	if len(cursor.sorters) > 0 {
		sorters := lo.Map(cursor.sorters, func(sorter record.Sorter, _ int) oas.RecordSorter {
			return oas.RecordSorter{
				Key:      sorter.Key,
				Ordering: string(sorter.Ordering),
			}
		})
		infoCursor.Sort = &sorters
	}

	if len(cursor.masks) > 0 {
		masks := lo.Map(cursor.masks, func(mask record.Mask, _ int) oas.RecordMask {
			return oas.RecordMask{
				Operation:       mask.Operation,
				PropertyPattern: mask.PropertyPattern,
			}
		})
		infoCursor.Masks = &masks
	}

	return infoCursor
}

type iter struct {
	cursor recordCursor

	index  int
	data   []record.Record
	loaded bool
	closed bool
}

func (iter *iter) Next(ctx context.Context) (record.R, error) {
	if iter.closed {
		return nil, errors.New("iterator is closed")
	}

	if err := iter.load(ctx); err != nil {
		return nil, fmt.Errorf("failed to load data: %w", err)
	}

	if iter.index >= len(iter.data) {
		return nil, record.ErrIterEmpty
	}

	rec := iter.data[iter.index]
	iter.index++

	return rec, nil
}

func (iter *iter) Close(ctx context.Context) error {
	iter.closed = true
	return nil
}

func (iter *iter) load(ctx context.Context) error {
	if iter.loaded {
		return nil
	}
	res, err := iter.cursor.httpClient.PostTasksInfoLoadWithResponse(ctx,
		oas.PostTasksInfoLoadJSONRequestBody{
			Cursor: iter.cursor.makeCursor(),
			Filter: iter.cursor.makeFilter(),
		})

	if err != nil {
		return fmt.Errorf("failed to make an http request: %w", err)
	}

	code := res.StatusCode()

	switch code {
	case http.StatusOK:
		iter.loaded = true
		iter.data = lo.Map(*res.JSON200, func(t oas.TaskInfo, _ int) record.R {
			return t
		})
		return nil
	case http.StatusInternalServerError:
		return fmt.Errorf("server failure: %s", res.JSON500.Message)
	default:
		return fmt.Errorf("unexpected response code: %d", code)
	}
}

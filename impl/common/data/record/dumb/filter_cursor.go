package dumb

import (
	"context"
	"github.com/samber/lo"
	"kantoku/common/data/record"
)

type FilterCursor struct {
	set     Set
	sorters []record.Sorter
	masks   []record.Mask
	skip    int
	limit   int
}

func (d FilterCursor) Skip(i int) record.Cursor[record.Record] {
	d.skip += i
	return d
}

func (d FilterCursor) Limit(i int) record.Cursor[record.Record] {
	d.limit = i
	return d
}

func (d FilterCursor) Mask(masks ...record.Mask) record.Cursor[record.Record] {
	d.masks = append(d.masks, masks...)
	return d
}

func (d FilterCursor) Sort(sorters ...record.Sorter) record.Cursor[record.Record] {
	d.sorters = sorters
	return d
}

func (d FilterCursor) Iter() record.Iter[record.Record] {
	return &FilterIter{
		cursor: d,
	}
}

func (d FilterCursor) Count(ctx context.Context) (int, error) {
	return len(d.eval()), nil
}

func (d FilterCursor) eval() []record.Record {
	data := d.set.eval()

	sorted(data, d.sorters)

	masked := lo.Map(data, func(r record.R, _ int) record.R {
		return mask(r, d.masks)
	})

	lim := d.limit
	if lim <= 0 {
		lim = len(masked)
	}

	return lo.Slice(masked, d.skip, d.skip+lim)
}

type FilterIter struct {
	index   int
	matched []record.R
	cursor  FilterCursor
}

func (d *FilterIter) Next(_ context.Context) (record.R, error) {
	if d.matched == nil {
		d.matched = d.cursor.eval()
	}

	if d.index >= len(d.matched) {
		return nil, record.ErrIterEmpty
	}

	res := d.matched[d.index]
	d.index++
	return res, nil
}

func (d *FilterIter) Close(ctx context.Context) error {
	return nil
}

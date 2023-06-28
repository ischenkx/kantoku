package dumb

import (
	"context"
	"github.com/samber/lo"
	"kantoku/common/data/record"
)

type DistinctCursor struct {
	set     Set
	keys    []string
	sorters []record.Sorter
	masks   []record.Mask
	skip    int
	limit   int
}

func (d DistinctCursor) Skip(i int) record.Cursor[record.Record] {
	d.skip += i
	return d
}

func (d DistinctCursor) Limit(i int) record.Cursor[record.Record] {
	d.limit = i
	return d
}

func (d DistinctCursor) Mask(masks ...record.Mask) record.Cursor[record.Record] {
	d.masks = append(d.masks, masks...)
	return d
}

func (d DistinctCursor) Sort(sorters ...record.Sorter) record.Cursor[record.Record] {
	d.sorters = sorters
	return d
}

func (d DistinctCursor) Iter() record.Iter[record.Record] {
	return &DistinctIter{
		cursor: d,
	}
}

func (d DistinctCursor) Count(ctx context.Context) (int, error) {
	return len(d.eval()), nil
}

func (d DistinctCursor) eval() []record.Record {
	data := d.set.eval()

	trie := newTrie(nil)
	distinct := lo.Filter(data, func(r record.Record, _ int) bool {
		values := lo.Map(keyMask(r, d.keys), func(entry record.E, _ int) any { return entry.Value })
		if trie.Exists(values) {
			return false
		}
		trie.Insert(values)

		return true
	})

	distinct = lo.Map(distinct, func(r record.R, _ int) record.R {
		return lo.SliceToMap(keyMask(r, d.keys), func(entry record.E) (string, any) {
			return entry.Name, entry.Value
		})
	})

	sorted(distinct, d.sorters)

	masked := lo.Map(distinct, func(r record.R, _ int) record.R {
		return mask(r, d.masks)
	})

	lim := d.limit
	if lim <= 0 {
		lim = len(masked)
	}

	return lo.Slice(masked, d.skip, d.skip+lim)
}

type DistinctIter struct {
	index   int
	matched []record.R
	cursor  DistinctCursor
}

func (d *DistinctIter) Next(_ context.Context) (record.R, error) {
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

func (d *DistinctIter) Close(ctx context.Context) error {
	return nil
}

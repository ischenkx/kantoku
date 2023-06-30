package dumb

import (
	"context"
	"github.com/samber/lo"
	"kantoku/common/data/record"
)

type Cursor struct {
	set     Set
	keys    []string
	sorters []record.Sorter
	masks   []record.Mask
	skip    int
	limit   int
}

func (d Cursor) Skip(i int) record.Cursor[record.Record] {
	d.skip += i
	return d
}

func (d Cursor) Limit(i int) record.Cursor[record.Record] {
	d.limit = i
	return d
}

func (d Cursor) Mask(masks ...record.Mask) record.Cursor[record.Record] {
	d.masks = append(d.masks, masks...)
	return d
}

func (d Cursor) Sort(sorters ...record.Sorter) record.Cursor[record.Record] {
	d.sorters = sorters
	return d
}

func (d Cursor) Iter() record.Iter[record.Record] {
	return &Iter{
		cursor: d,
	}
}

func (d Cursor) Count(ctx context.Context) (int, error) {
	return len(d.eval()), nil
}

func (d Cursor) eval() []record.Record {
	data := d.set.eval()

	if d.keys != nil {
		trie := newTrie(nil)
		data = lo.Filter(data, func(r record.Record, _ int) bool {
			values := lo.Map(keyMask(r, d.keys), func(entry record.E, _ int) any { return entry.Value })
			if trie.Exists(values) {
				return false
			}
			trie.Insert(values)

			return true
		})

		data = lo.Map(data, func(r record.R, _ int) record.R {
			return lo.SliceToMap(keyMask(r, d.keys), func(entry record.E) (string, any) {
				return entry.Name, entry.Value
			})
		})
	}

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

type Iter struct {
	index   int
	matched []record.R
	cursor  Cursor
}

func (d *Iter) Next(_ context.Context) (record.R, error) {
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

func (d *Iter) Close(ctx context.Context) error {
	return nil
}

package dumb

import (
	"context"
	"fmt"
	"github.com/ischenkx/kantoku/pkg/common/data/codec"
	"github.com/ischenkx/kantoku/pkg/common/data/record"
	"github.com/samber/lo"
)

var _ record.Cursor[int] = Cursor[int]{}

type Cursor[Item any] struct {
	set     Set[Item]
	keys    []string
	sorters []record.Sorter
	masks   []record.Mask
	Codec   codec.Codec[Item, record.R]
	skip    int
	limit   int
}

func (d Cursor[Item]) Skip(i int) record.Cursor[Item] {
	d.skip += i
	return d
}

func (d Cursor[Item]) Limit(i int) record.Cursor[Item] {
	d.limit = i
	return d
}

func (d Cursor[Item]) Mask(masks ...record.Mask) record.Cursor[Item] {
	d.masks = append(d.masks, masks...)
	return d
}

func (d Cursor[Item]) Sort(sorters ...record.Sorter) record.Cursor[Item] {
	d.sorters = sorters
	return d
}

func (d Cursor[Item]) Iter() record.Iter[Item] {
	return &Iter[Item]{
		cursor: d,
	}
}

func (d Cursor[Item]) Count(ctx context.Context) (int, error) {
	res, err := d.eval()
	if err != nil {
		return 0, err
	}

	return len(res), nil
}

func (d Cursor[Item]) eval() ([]Item, error) {
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

	var result []Item
	for _, encoded := range masked {
		item, err := d.Codec.Decode(encoded)
		if err != nil {
			return nil, fmt.Errorf("failed to decode: %w", err)
		}
		result = append(result, item)
	}

	return lo.Slice(result, d.skip, d.skip+lim), nil
}

type Iter[Item any] struct {
	index   int
	matched []Item
	cursor  Cursor[Item]
}

func (d *Iter[Item]) Next(_ context.Context) (Item, error) {
	if d.matched == nil {
		matched, err := d.cursor.eval()
		if err != nil {
			var zero Item
			return zero, fmt.Errorf("failed to eval: %w", err)
		}
		d.matched = matched
	}

	if d.index >= len(d.matched) {
		var zero Item
		return zero, record.ErrIterEmpty
	}

	res := d.matched[d.index]
	d.index++
	return res, nil
}

func (d *Iter[Item]) Close(ctx context.Context) error {
	return nil
}

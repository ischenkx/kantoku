package record

import "context"

type Cursor[Item any] interface {
	Skip(int) Cursor[Item]
	Limit(int) Cursor[Item]
	Mask(masks ...Mask) Cursor[Item]
	Sort(sorters ...Sorter) Cursor[Item]
	Iter() Iter[Item]
	Count(ctx context.Context) (int, error)
}

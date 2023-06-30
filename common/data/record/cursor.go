package record

import "context"

type Cursor[Item any] interface {
	Skip(int) Cursor[Item]
	Limit(int) Cursor[Item]
	// Mask is a method that allows dynamically include/exclude some entries by their name
	//
	// NOTE: Following operations are equivalent:
	//     cursor.Mask(masks1...).Mass(masks2...)
	//     cursor.Mask(append(masks1, masks2...)...)
	Mask(masks ...Mask) Cursor[Item]
	Sort(sorters ...Sorter) Cursor[Item]
	Iter() Iter[Item]
	Count(ctx context.Context) (int, error)
}

package record

import "context"

type Cursor[Item any] interface {
	Skip(int) Cursor[Item]
	Limit(int) Cursor[Item]
	// Mask is a method that allows dynamically include/exclude some entries by their name
	// There are two modes: exclude and include.
	// If mask is empty, all fields will be present in result of query.
	// If you apply excluding mask, fields you masked will be ignored.
	// Similarly, if you include some fields, only they will be present.
	// You cannot apply both mask types simultaneously.
	//
	// Applying multiple masks works as logical OR, so following operations are equivalent:
	//     cursor.Mask(masks1...).Mass(masks2...)
	//     cursor.Mask(append(masks1, masks2...)...)
	Mask(masks ...Mask) Cursor[Item]
	Sort(sorters ...Sorter) Cursor[Item]
	Iter() Iter[Item]
	Count(ctx context.Context) (int, error)
}

package record

import "context"

type Set[Item any] interface {
	Filter(rec Record) Set[Item]
	Erase(ctx context.Context) error
	// Update updates or inserts filtered values.
	//
	// If upsert is not nil and no records are matched then a new value is inserted (upsert)
	// and then updated (update)
	Update(ctx context.Context, update, upsert R) error
	Distinct(key ...string) Cursor[Item]
	Cursor() Cursor[Item]
}

package record

import "context"

type Set interface {
	Filter(...Entry) Set
	Erase(ctx context.Context) error
	// Update updates or inserts filtered values.
	//
	// If upsert is not nil and no records are matched then a new value is inserted (upsert)
	// and updated (update)
	Update(ctx context.Context, update, upsert R) error
	Distinct(key ...string) Cursor[Record]
	Cursor() Cursor[Record]
}

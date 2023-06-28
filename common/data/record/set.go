package record

import "context"

type Set interface {
	Filter(...Entry) Set
	Erase(ctx context.Context) error
	Update(ctx context.Context, update, upsert R) error
	Distinct(key ...string) Cursor[Record]
	Cursor() Cursor[Record]
}

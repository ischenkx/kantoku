package cellkv

import (
	"context"
	"kantoku/common/data/kv"
	"kantoku/common/util"
	"kantoku/core/cell"
)

type DB[T any] struct {
	cells   cell.Storage[T]
	id2cell kv.Database[string]
}

func (db *DB[T]) Set(ctx context.Context, id string, item T) (T, error) {
	cell, err := db.cells.Make(ctx, item)
	if err != nil {
		return util.Default[T](), err
	}
	if _, err := db.id2cell.Set(ctx, id, cell); err != nil {
		return util.Default[T](), err
	}

	return item, nil
}

func (db *DB[T]) Del(ctx context.Context, id string) error {
	return db.id2cell.Del(ctx, id)
}

func (db *DB[T]) Get(ctx context.Context, id string) (T, error) {
	cellID, err := db.id2cell.Get(ctx, id)
	if err != nil {
		return util.Default[T](), err
	}
	cell, err := db.cells.Get(ctx, cellID)
	if err != nil {
		return util.Default[T](), err
	}
	return cell.Data, nil
}

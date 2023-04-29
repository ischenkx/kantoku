package cellkv

import (
	"context"
	"kantoku/common/data"
	"kantoku/common/data/cell"
	"kantoku/common/data/kv"
	"kantoku/common/util"
)

type DB[T any] struct {
	cells   cell.Storage[T]
	id2cell kv.Database[string, string]
}

func (db *DB[T]) Set(ctx context.Context, id string, item T) error {
	cellID, err := db.cells.Make(ctx, item)
	if err != nil {
		return err
	}
	if err := db.id2cell.Set(ctx, id, cellID); err != nil {
		return err
	}

	return nil
}

func (db *DB[T]) GetOrSet(ctx context.Context, id string, item T) (T, error) {
	cellID, err := db.id2cell.Get(ctx, id)

	if err == nil {
		val, err := db.cells.Get(ctx, cellID)
		return val.Data, err
	} else if err != data.NotFoundErr {
		return util.Default[T](), err
	}

	cellID, err = db.cells.Make(ctx, item)
	if err != nil {
		return util.Default[T](), err
	}

	cellID, err = db.id2cell.GetOrSet(ctx, id, cellID)
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
	c, err := db.cells.Get(ctx, cellID)
	if err != nil {
		return util.Default[T](), err
	}
	return c.Data, nil
}

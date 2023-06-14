package redikv

import (
	"context"
	"github.com/redis/go-redis/v9"
	"kantoku/common/codec"
	"kantoku/common/data"
	"kantoku/common/data/kv"
	"kantoku/common/util"
)

type DB[T any] struct {
	client  redis.UniversalClient
	codec   codec.Codec[T, []byte]
	setName string
}

var _ kv.Database[string, int] = (*DB[int])(nil)

func New[T any](client redis.UniversalClient, codec codec.Codec[T, []byte], setName string) *DB[T] {
	return &DB[T]{
		client:  client,
		codec:   codec,
		setName: setName,
	}
}

func (db *DB[T]) Set(ctx context.Context, id string, item T) error {
	data, err := db.codec.Encode(item)
	if err != nil {
		return err
	}
	if cmd := db.client.HSet(ctx, db.setName, id, data); cmd.Err() != nil {
		return cmd.Err()
	}
	return nil
}

func (db *DB[T]) GetOrSet(ctx context.Context, id string, item T) (T, bool, error) {
	val, err := db.codec.Encode(item)
	if err != nil {
		return util.Default[T](), false, err
	}

	cmd := db.client.HSetNX(ctx, db.setName, id, val)
	if cmd.Err() != nil {
		return util.Default[T](), false, cmd.Err()
	}

	result, err := db.Get(ctx, id)
	return result, cmd.Val(), err
}

func (db *DB[T]) Get(ctx context.Context, id string) (T, error) {
	cmd := db.client.HGet(ctx, db.setName, id)
	if cmd.Err() != nil {
		err := cmd.Err()
		if err == redis.Nil {
			err = data.NotFoundErr
		}
		return util.Default[T](), err
	}
	raw, err := cmd.Bytes()
	if cmd.Err() != nil {
		return util.Default[T](), err
	}
	val, err := db.codec.Decode(raw)
	if err != nil {
		return util.Default[T](), err
	}
	return val, nil
}

func (db *DB[T]) Del(ctx context.Context, id string) error {
	cmd := db.client.HDel(ctx, db.setName, id)
	return cmd.Err()
}

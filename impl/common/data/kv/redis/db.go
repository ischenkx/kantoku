package redikv

import (
	"bytes"
	"context"
	"github.com/redis/go-redis/v9"
	"kantoku/common/codec"
	"kantoku/common/util"
)

type DB[T any] struct {
	client  redis.UniversalClient
	codec   codec.Codec[T]
	setName string
}

func New[T any](client redis.UniversalClient, codec codec.Codec[T], setName string) *DB[T] {
	return &DB[T]{
		client:  client,
		codec:   codec,
		setName: setName,
	}
}

func (db *DB[T]) Set(ctx context.Context, id string, item T) (T, error) {
	data, err := db.codec.Encode(item)
	if err != nil {
		return util.Default[T](), err
	}
	if cmd := db.client.HSet(ctx, db.setName, id, data); cmd.Err() != nil {
		return util.Default[T](), cmd.Err()
	}
	return item, nil
}

func (db *DB[T]) Get(ctx context.Context, id string) (T, error) {
	cmd := db.client.HGet(ctx, db.setName, id)
	if cmd.Err() != nil {
		return util.Default[T](), cmd.Err()
	}
	raw, err := cmd.Bytes()
	if cmd.Err() != nil {
		return util.Default[T](), err
	}
	data, err := db.codec.Decode(bytes.NewReader(raw))
	if err != nil {
		return util.Default[T](), err
	}
	return data, nil
}

func (db *DB[T]) Del(ctx context.Context, id string) error {
	cmd := db.client.HDel(ctx, db.setName, id)
	return cmd.Err()
}

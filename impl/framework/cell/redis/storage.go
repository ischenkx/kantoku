package redicell

import (
	"bytes"
	"context"
	"github.com/redis/go-redis/v9"
	"kantoku/common/codec"
	"kantoku/common/util"
	"kantoku/framework/cell"
)

type Storage[T any] struct {
	client redis.UniversalClient
	codec  codec.Codec[T]
}

func NewStorage[T any](client redis.UniversalClient, codec codec.Codec[T]) *Storage[T] {
	return &Storage[T]{
		client: client,
		codec:  codec,
	}
}

func (s *Storage[T]) Set(ctx context.Context, id string, data T) error {
	encoded, err := s.codec.Encode(data)
	if err != nil {
		return err
	}

	err = s.client.Set(ctx, id, encoded, 0).Err()
	if err != nil {
		return err
	}
	return nil
}

func (s *Storage[T]) Get(ctx context.Context, id string) (cell.Cell[T], error) {
	encoded, err := s.client.Get(ctx, id).Bytes()
	if err != nil {
		return util.Default[cell.Cell[T]](), err
	}

	data, err := s.codec.Decode(bytes.NewReader(encoded))
	if err != nil {
		return util.Default[cell.Cell[T]](), err
	}

	return cell.Cell[T]{
		ID:   id,
		Data: data,
	}, nil
}

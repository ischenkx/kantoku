package redicell

import (
	"context"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	"kantoku/common/codec"
	"kantoku/common/data/cell"
	"kantoku/common/util"
)

type Storage[T any] struct {
	client redis.UniversalClient
	codec  codec.Codec[T, []byte]
}

func New[T any](client redis.UniversalClient, codec codec.Codec[T, []byte]) *Storage[T] {
	return &Storage[T]{
		client: client,
		codec:  codec,
	}
}

func (s *Storage[T]) Make(ctx context.Context, data T) (string, error) {
	id := uuid.New().String()

	encoded, err := s.codec.Encode(data)
	if err != nil {
		return "", err
	}

	err = s.client.Set(ctx, id, encoded, 0).Err()
	if err != nil {
		return id, err
	}
	return id, nil
}

func (s *Storage[T]) Get(ctx context.Context, id string) (cell.Cell[T], error) {
	encoded, err := s.client.Get(ctx, id).Bytes()
	if err != nil {
		return util.Default[cell.Cell[T]](), err
	}

	data, err := s.codec.Decode(encoded)
	if err != nil {
		return util.Default[cell.Cell[T]](), err
	}

	return cell.Cell[T]{
		ID:   id,
		Data: data,
	}, nil
}

package redis

import (
	"context"
	"github.com/redis/go-redis/v9"
	"github.com/satori/go.uuid"
	"hayaku/l0/cell"
)

type Storage struct {
	client redis.UniversalClient
}

func NewStorage(client redis.UniversalClient) *Storage {
	return &Storage{
		client: client,
	}
}

func (s *Storage) Create(ctx context.Context, data []byte) (string, error) {
	id := uuid.NewV4().String()
	err := s.client.Set(ctx, id, data, 0).Err()
	if err != nil {
		return "", err
	}
	return id, nil
}

func (s *Storage) Get(ctx context.Context, id string) (*cell.Cell, error) {
	data, err := s.client.Get(ctx, id).Bytes()
	if err != nil {
		return nil, err
	}
	return &cell.Cell{
		ID:   id,
		Data: data,
	}, nil
}

func (s *Storage) Set(ctx context.Context, cell *cell.Cell) error {
	return s.client.Set(ctx, cell.ID, cell.Data, 0).Err()
}

func (s *Storage) Delete(ctx context.Context, id string) error {
	return s.client.Del(ctx, id).Err()
}

func (s *Storage) Close() error {
	return s.client.Close()
}

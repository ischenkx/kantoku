package specification

import (
	"context"
	"errors"
	"fmt"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PostgresBinaryStorage struct {
	DB    *pgxpool.Pool
	Table string
	_     struct{}
}

func (s *PostgresBinaryStorage) Get(ctx context.Context, id string) ([]byte, error) {
	var data []byte
	query := fmt.Sprintf("SELECT data FROM %s WHERE id = $1", s.Table)
	err := s.DB.QueryRow(ctx, query, id).Scan(&data)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	}
	return data, err
}

func (s *PostgresBinaryStorage) GetAll(ctx context.Context) ([][]byte, error) {
	query := fmt.Sprintf("SELECT data FROM %s", s.Table)
	rows, err := s.DB.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result [][]byte
	for rows.Next() {
		var data []byte
		if err := rows.Scan(&data); err != nil {
			return nil, err
		}
		result = append(result, data)
	}
	return result, nil
}

func (s *PostgresBinaryStorage) Add(ctx context.Context, id string, data []byte) error {
	query := fmt.Sprintf("INSERT INTO %s (id, data) VALUES ($1, $2)", s.Table)
	_, err := s.DB.Exec(ctx, query, id, data)
	return err
}

func (s *PostgresBinaryStorage) Remove(ctx context.Context, id string) error {
	query := fmt.Sprintf("DELETE FROM %s WHERE id = $1", s.Table)
	_, err := s.DB.Exec(ctx, query, id)
	return err
}

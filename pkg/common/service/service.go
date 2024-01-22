package service

import (
	"context"
	"log/slog"
)

type Service interface {
	Run(ctx context.Context) error
	Logger() *slog.Logger
	ID() string
	Name() string
}

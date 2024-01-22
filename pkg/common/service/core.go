package service

import (
	"context"
	"fmt"
	"log/slog"
)

type Core struct {
	logger *slog.Logger
	id     string
	name   string
}

func NewCore(name string, id string, logger *slog.Logger) Core {
	return Core{
		logger: logger,
		id:     id,
		name:   name,
	}
}

func (core Core) Run(ctx context.Context) error {
	return fmt.Errorf("not implemented")
}

func (core Core) ID() string {
	return core.id
}

func (core Core) Name() string {
	return core.name
}

func (core Core) Logger() *slog.Logger {
	return core.logger
}

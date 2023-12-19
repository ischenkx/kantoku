package alpha

import (
	"context"
	"fmt"
	"kantoku/common/data/record"
)

type Alpha struct {
	id      string
	manager *Manager
}

func (alpha Alpha) ID() string {
	return alpha.id
}

func (alpha Alpha) Data(ctx context.Context) ([]byte, error) {
	data, err := alpha.manager.storage.Get(ctx, alpha.ID())
	if err != nil {
		return nil, fmt.Errorf("failed to load the object: %s", err)
	}

	return data, nil
}

func (alpha Alpha) Info() record.Set {
	return alpha.manager.Info().Filter(record.E{TaskIDProperty, alpha.ID()})
}

func (alpha Alpha) Manager() *Manager {
	return alpha.manager
}

func (alpha Alpha) Result(ctx context.Context) (Result, error) {
	return alpha.manager.results.Get(ctx, alpha.ID())
}

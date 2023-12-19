package gamma

import (
	"context"
	"kantoku/common/data/deps"
	"kantoku/core/beta"
)

type Gamma struct {
	id      string
	manager *Manager
}

func (gamma Gamma) ID() string {
	return gamma.id
}

func (gamma Gamma) Beta() beta.Beta {
	return gamma.manager.betas.Get(gamma.id)
}

func (gamma Gamma) Dependencies(ctx context.Context) []deps.Dependency {

}

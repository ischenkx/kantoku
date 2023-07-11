package demon

import (
	"context"
)

type Manager struct {
	demons       []Demon
	orchestrator Deployer
}

func NewManager(orchestrator Deployer) *Manager {
	return &Manager{orchestrator: orchestrator}
}

func (manager *Manager) Register(demons ...Demon) *Manager {
	for _, demon := range demons {
		manager.register(demon)
	}
	return manager
}

func (manager *Manager) List() []Demon {
	return manager.demons
}

func (manager *Manager) Run(ctx context.Context) error {
	return manager.orchestrator.Deploy(ctx, manager.demons...)
}

func (manager *Manager) register(demon Demon) {
	for _, other := range manager.List() {
		if other.Eq(demon) {
			return
		}
	}
	manager.demons = append(manager.demons, demon)
}

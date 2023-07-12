package demon

import (
	"context"
)

type Manager struct {
	demons       []Demon
	deployer Deployer
}

func NewManager(deployer Deployer) *Manager {
	return &Manager{deployer: deployer}
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
	return manager.deployer.Deploy(ctx, manager.demons...)
}

func (manager *Manager) register(demon Demon) {
	for _, other := range manager.List() {
		if other.Eq(demon) {
			return
		}
	}
	manager.demons = append(manager.demons, demon)
}

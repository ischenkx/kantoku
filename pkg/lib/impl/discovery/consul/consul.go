package consul

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/hashicorp/consul/api"
	"github.com/ischenkx/kantoku/pkg/lib/discovery"
	"sync"
)

type Hub struct {
	Consul *api.Client
	mu     sync.Mutex
}

func (hub *Hub) Register(ctx context.Context, info discovery.ServiceInfo) error {
	hub.mu.Lock()
	defer hub.mu.Unlock()

	agent := hub.Consul.Agent()

	check := &api.AgentServiceCheck{
		CheckID:                        "health-check-" + info.ID,
		TTL:                            "15s",
		DeregisterCriticalServiceAfter: "1m",
	}

	encodedInfo, err := json.Marshal(info.Info)
	if err != nil {
		return fmt.Errorf("failed to encode info: %w", err)
	}

	encodedStatus, err := json.Marshal(info.Status)
	if err != nil {
		return fmt.Errorf("failed to encode status: %w", err)
	}

	meta := map[string]string{
		"info":   string(encodedInfo),
		"status": string(encodedStatus),
	}

	serviceRegistration := &api.AgentServiceRegistration{
		ID:    info.ID,
		Name:  info.Name,
		Tags:  []string{"service"},
		Check: check,
		Meta:  meta,
	}

	err = agent.ServiceRegister(serviceRegistration)
	if err != nil {
		return err
	}

	fmt.Printf("Service registered: %+v\n", info)
	return nil
}

func (hub *Hub) Load(ctx context.Context) ([]discovery.ServiceInfo, error) {
	return nil, fmt.Errorf("not implemented")
}

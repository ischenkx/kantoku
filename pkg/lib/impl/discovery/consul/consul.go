package consul

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/hashicorp/consul/api"
	"github.com/ischenkx/kantoku/pkg/lib/discovery"
	"github.com/mitchellh/mapstructure"
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
		Status:                         "passing",
		TTL:                            "15s",
		DeregisterCriticalServiceAfter: "30s",
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

	addr, port, ok := parseLocation(info.Info)
	if ok {
		serviceRegistration.Address = addr
		serviceRegistration.Port = port
	}

	err = agent.ServiceRegister(serviceRegistration)
	if err != nil {
		return err
	}

	//fmt.Printf("Service registered: %+v\n", info)
	return nil
}

func (hub *Hub) Load(ctx context.Context) ([]discovery.ServiceInfo, error) {
	return nil, fmt.Errorf("not implemented")
}

func parseLocation(info map[string]any) (addr string, port int, ok bool) {
	var cfg struct {
		Addr string
		Port int
	}
	if err := mapstructure.Decode(info, &cfg); err != nil {
		return "", 0, false
	}

	return cfg.Addr, cfg.Port, true
}

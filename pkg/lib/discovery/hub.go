package discovery

import "context"

type ServiceInfo struct {
	ID     string
	Name   string
	Info   map[string]any
	Status map[string]any
}

type Hub interface {
	Register(ctx context.Context, info ServiceInfo) error
	Load(ctx context.Context) ([]ServiceInfo, error)
}

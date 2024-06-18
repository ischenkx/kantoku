package mock

import (
	"context"
	"fmt"
	"github.com/ischenkx/kantoku/pkg/core/resource"
	"github.com/samber/lo"
)

type Storage struct {
	idCounter int
	data      map[resource.ID]resource.Resource
}

func NewMockStorage() *Storage {
	return &Storage{
		idCounter: 0,
		data:      map[resource.ID]resource.Resource{},
	}
}

func (s *Storage) Load(ctx context.Context, ids ...resource.ID) ([]resource.Resource, error) {
	return lo.Map(ids, func(id resource.ID, _ int) resource.Resource {
		res, has := s.data[id]
		if has {
			return res
		} else {
			return resource.Resource{
				Data:   nil,
				ID:     id,
				Status: resource.DoesNotExist,
			}
		}
	}), nil
}

func (s *Storage) Alloc(ctx context.Context, amount int) ([]resource.ID, error) {
	ids := make([]resource.ID, amount)
	for i := 0; i < amount; i++ {
		s.idCounter++
		ids[i] = fmt.Sprint("resource-", s.idCounter)
	}
	return ids, nil
}

func (s *Storage) Init(ctx context.Context, resources []resource.Resource) error {
	for _, res := range resources {
		res.Status = resource.Ready
		s.data[res.ID] = res
	}
	return nil
}

func (s *Storage) Dealloc(ctx context.Context, ids []resource.ID) error {
	for _, id := range ids {
		delete(s.data, id)
	}
	return nil
}

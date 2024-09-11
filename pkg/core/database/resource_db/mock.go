package resourcedb

import (
	"context"
	"fmt"
	"github.com/ischenkx/kantoku/pkg/core"
	"github.com/samber/lo"
)

type MockDB struct {
	idCounter int
	data      map[string]core.Resource
}

func NewMockDB() *MockDB {
	return &MockDB{
		idCounter: 0,
		data:      map[string]core.Resource{},
	}
}

func (s *MockDB) Load(ctx context.Context, ids ...string) ([]core.Resource, error) {
	return lo.Map(ids, func(id string, _ int) core.Resource {
		res, has := s.data[id]
		if has {
			return res
		} else {
			return core.Resource{
				Data:   nil,
				ID:     id,
				Status: core.ResourceStatuses.DoesNotExist,
			}
		}
	}), nil
}

func (s *MockDB) Alloc(ctx context.Context, amount int) ([]string, error) {
	ids := make([]string, amount)
	for i := 0; i < amount; i++ {
		s.idCounter++
		ids[i] = fmt.Sprint("resource_db-", s.idCounter)
	}
	return ids, nil
}

func (s *MockDB) Init(ctx context.Context, resources []core.Resource) error {
	for _, res := range resources {
		res.Status = core.ResourceStatuses.Ready
		s.data[res.ID] = res
	}
	return nil
}

func (s *MockDB) Dealloc(ctx context.Context, ids []string) error {
	for _, id := range ids {
		delete(s.data, id)
	}
	return nil
}

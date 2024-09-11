package future

import (
	"context"
	"github.com/ischenkx/kantoku/pkg/common/data/codec"
	"github.com/ischenkx/kantoku/pkg/core"
)

// Storage manages future-resource_db mapping. It is not thread safe
type Storage struct {
	id2future   map[fid]AbstractFuture
	id2resource map[fid]*core.Resource
	isSaved     map[fid]bool

	assignedLog []string
}

func NewStorage() Storage {
	return Storage{
		id2future:   map[fid]AbstractFuture{},
		id2resource: map[fid]*core.Resource{},
		isSaved:     map[fid]bool{},
		assignedLog: []string{},
	}
}

func (s *Storage) AddFuture(fut AbstractFuture) {
	s.id2future[fut.getId()] = fut
}

func (s *Storage) AssignResource(fut AbstractFuture, res *core.Resource, saved bool) {
	s.id2resource[fut.getId()] = res
	s.isSaved[fut.getId()] = saved
}

// Allocate ids for resources without them, empty resources are assigned to futures without them
func (s *Storage) Allocate(ctx context.Context, storage core.ResourceDB) error {
	for id, _ := range s.id2future {
		res := s.id2resource[id]
		if res == nil {
			s.id2resource[id] = &core.Resource{}
		}
	}

	unallocated := make([]*core.Resource, 0)
	for _, res := range s.id2resource {
		if res.ID == "" {
			unallocated = append(unallocated, res)
		}
	}
	ids, err := storage.Alloc(ctx, len(unallocated))
	s.assignedLog = append(s.assignedLog, ids...)
	if err != nil {
		return err
	}
	for i := 0; i < len(ids); i++ {
		unallocated[i].Status = core.ResourceStatuses.Allocated
		unallocated[i].ID = ids[i]
	}
	return nil
}

// Encode all filled futures. It will create resources, or fill Data field for existing ones.
// Not filled futures and resources with Data are skipped.
func (s *Storage) Encode(codec codec.Codec[any, []byte]) error {
	for id, fut := range s.id2future {
		res := s.id2resource[id]
		if res == nil {
			res = &core.Resource{}
		}

		if !fut.IsFilled() || res.Data != nil {
			continue
		}

		data, err := fut.Encode(codec)
		if err != nil {
			return err
		}
		res.Data = data

		s.id2resource[id] = res
	}
	return nil
}

// Save all resources that are not marked as saved and have data
func (s *Storage) Save(ctx context.Context, storage core.ResourceDB) error {
	toSave := make([]core.Resource, 0)
	for id, res := range s.id2resource {
		if s.id2resource[id].ID != "" && !s.isSaved[id] && res.Data != nil {
			toSave = append(toSave, *res)
		}
	}
	err := storage.Init(ctx, toSave)
	if err != nil {
		return err
	}
	for id := range s.id2resource {
		if s.id2resource[id].ID != "" && !s.isSaved[id] {
			s.isSaved[id] = true
		}
	}
	return nil
}

func (s *Storage) GetResource(fut AbstractFuture) *core.Resource {
	return s.id2resource[fut.getId()]
}

func (s *Storage) HasFuture(fut AbstractFuture) bool {
	_, has := s.id2future[fut.getId()]
	return has
}

func (s *Storage) Rollback(ctx context.Context, storage core.ResourceDB) error {
	err := storage.Dealloc(ctx, s.assignedLog)
	if err != nil {
		return err
	}
	s.assignedLog = []string{}
	return nil
}

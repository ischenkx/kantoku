package future

import (
	"context"
	"github.com/ischenkx/kantoku/pkg/system/kernel/resource"
)

// Storage manages future-resource mapping. It is not thread safe
type Storage struct {
	id2future   map[fid]Future[any]
	id2resource map[fid]*resource.Resource
	isSaved     map[fid]bool

	assignedLog []resource.ID
}

func NewStorage() Storage {
	return Storage{
		id2future:   map[fid]Future[any]{},
		id2resource: map[fid]*resource.Resource{},
		isSaved:     map[fid]bool{},
		assignedLog: []resource.ID{},
	}
}

func (s Storage) AddFuture(fut Future[any]) {
	s.id2future[fut.id] = fut
}

func (s Storage) AssignResource(fut Future[any], res *resource.Resource, saved bool) {
	s.id2resource[fut.id] = res
	s.isSaved[fut.id] = saved
}

// Allocate ids for resources without them
func (s Storage) Allocate(ctx context.Context, storage resource.Storage) error {
	unallocated := make([]*resource.Resource, 0)
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
		unallocated[i].ID = ids[i]
	}
	return nil
}

// Save all resources that are not marked as saved
func (s Storage) Save(ctx context.Context, storage resource.Storage) error {
	notSaved := make([]resource.Resource, 0)
	for id, res := range s.id2resource {
		if !s.isSaved[id] {
			notSaved = append(notSaved, *res)
		}
	}
	err := storage.Init(ctx, notSaved)
	if err != nil {
		return err
	}
	for id := range s.id2resource {
		if !s.isSaved[id] {
			s.isSaved[id] = true
		}
	}
	return nil
}

func (s Storage) GetResource(fut Future[any]) *resource.Resource {
	return s.id2resource[fut.id]
}

func (s Storage) HasFuture(fut Future[any]) bool {
	_, has := s.id2future[fut.id]
	return has
}

func (s Storage) Rollback(ctx context.Context, storage resource.Storage) error {
	err := storage.Dealloc(ctx, s.assignedLog)
	if err != nil {
		return err
	}
	s.assignedLog = []resource.ID{}
	return nil
}

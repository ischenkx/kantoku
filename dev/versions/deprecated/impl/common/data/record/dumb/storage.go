package dumb

import (
	"context"
	"github.com/samber/lo"
	"kantoku/common/data/record"
)

var _ record.Storage = (*Storage)(nil)

type Storage struct {
	data []record.R
}

func New() *Storage {
	return &Storage{}
}

func (s *Storage) List() []record.R {
	return s.data
}

func (s *Storage) Sample() record.R {
	return lo.Sample(s.data).Copy()
}

func (s *Storage) Insert(ctx context.Context, rec record.Record) error {
	if len(rec) == 0 {
		return nil
	}
	s.data = append(s.data, rec)
	return nil
}

func (s *Storage) Filter(entries ...record.Entry) record.Set {
	return newSet(s).Filter(entries...)
}

func (s *Storage) Erase(ctx context.Context) error {
	return newSet(s).Erase(ctx)
}

func (s *Storage) Update(ctx context.Context, update, upsert record.R) error {
	return newSet(s).Update(ctx, update, upsert)
}

func (s *Storage) Distinct(keys ...string) record.Cursor[record.Record] {
	return newSet(s).Distinct(keys...)
}

func (s *Storage) Cursor() record.Cursor[record.Record] {
	return newSet(s).Cursor()
}

func (s *Storage) each(f func(r record.Record, index int)) {
	lo.ForEach(s.data, f)
}

func (s *Storage) filter(f func(r record.Record, index int) bool) {
	s.data = lo.Filter(s.data, f)
}

func (s *Storage) update(f func(r record.Record, index int) record.Record) {
	s.data = lo.Map(s.data, f)
}

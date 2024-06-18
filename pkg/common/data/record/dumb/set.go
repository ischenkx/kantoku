package dumb

import (
	"context"
	"github.com/ischenkx/kantoku/pkg/common/data/record"
)

var _ record.Set[int] = Set[int]{}

type Set[Item any] struct {
	filters [][]record.Entry
	storage *Storage[Item]
}

func newSet[Item any](s *Storage[Item]) Set[Item] {
	return Set[Item]{storage: s}
}

func (set Set[Item]) Filter(rec record.Record) record.Set[Item] {
	for key, val := range rec {

	}
	set.filters = append(set.filters, entries)
	return set
}

func (set Set[Item]) Erase(_ context.Context) error {
	set.storage.filter(func(r record.Record, _ int) bool {
		return !matches(r, set.filters)
	})
	return nil
}

func (set Set[Item]) Update(ctx context.Context, update, upsert record.R) error {
	matched := false
	set.storage.update(func(r record.Record, _ int) record.Record {
		if !matches(r, set.filters) {
			return r
		}
		matched = true

		for key, value := range update {
			r[key] = value
		}
		return r
	})

	if !matched && upsert != nil {
		newRecord := record.R{}

		for key, value := range upsert {
			newRecord[key] = value
		}

		for key, value := range update {
			newRecord[key] = value
		}

		if err := set.storage.Insert(ctx, newRecord); err != nil {
			return err
		}
	}

	return nil
}

func (set Set[Item]) Distinct(keys ...string) record.Cursor[record.Record] {
	return Cursor{
		set:  set,
		keys: keys,
	}
}

func (set Set[Item]) Cursor() record.Cursor[record.Record] {
	return Cursor{
		set: set,
	}
}

func (set Set[Item]) eval() []record.R {
	var matched []record.R
	set.storage.each(func(r record.Record, _ int) {
		if matches(r, set.filters) {
			matched = append(matched, r.Copy())
		}
	})

	return matched
}

package dumb

import (
	"context"
	"github.com/ischenkx/kantoku/pkg/common/data/record"
)

var _ record.Set = Set{}

type Set struct {
	filters [][]record.Entry
	storage *Storage
}

func newSet(s *Storage) Set {
	return Set{storage: s}
}

func (set Set) Filter(entries ...record.Entry) record.Set {
	set.filters = append(set.filters, entries)
	return set
}

func (set Set) Erase(_ context.Context) error {
	set.storage.filter(func(r record.Record, _ int) bool {
		return !matches(r, set.filters)
	})
	return nil
}

func (set Set) Update(ctx context.Context, update, upsert record.R) error {
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

func (set Set) Distinct(keys ...string) record.Cursor[record.Record] {
	return Cursor{
		set:  set,
		keys: keys,
	}
}

func (set Set) Cursor() record.Cursor[record.Record] {
	return Cursor{
		set: set,
	}
}

func (set Set) eval() []record.R {
	var matched []record.R
	set.storage.each(func(r record.Record, _ int) {
		if matches(r, set.filters) {
			matched = append(matched, r.Copy())
		}
	})

	return matched
}

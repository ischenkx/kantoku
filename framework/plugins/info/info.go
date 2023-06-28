package info

import (
	"context"
	"fmt"
	"github.com/samber/lo"
	"kantoku/common/data/record"
)

type Dict = record.R
type Entry = record.E

type Info struct {
	storage *Storage
	id      string
}

func (info Info) Get(ctx context.Context, property string) (any, error) {
	item, err := info.Load(ctx, property)
	if err != nil {
		return nil, err
	}
	return item[property], nil
}

func (info Info) Erase(ctx context.Context) error {
	return info.set().Erase(ctx)
}

// Load returns a map that includes properties (if properties are empty then all properties are included)
func (info Info) Load(ctx context.Context, properties ...string) (Dict, error) {
	iter := info.set().
		Cursor().
		Mask(
			lo.Map(properties, func(prop string, _ int) record.Mask { return record.Include(prop) })...,
		).
		Iter()
	defer iter.Close(ctx)

	item, err := iter.Next(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get the value: %s", err)
	}

	for key, value := range item {
		if value == nil {
			delete(item, key)
		}
	}

	return Dict(item), nil
}

func (info Info) Set(ctx context.Context, entries ...Entry) error {
	update := lo.SliceToMap(entries, func(entry Entry) (string, any) {
		return entry.Name, entry.Value
	})
	upsert := record.R{info.storage.settings.IdProperty: info.id}

	err := info.set().Update(ctx, update, upsert)
	if err != nil {
		return fmt.Errorf("failed to update the record: %s", err)
	}

	return nil
}

// Del is equal to setting to nil
func (info Info) Del(ctx context.Context, properties ...string) error {
	if lo.Contains(properties, info.storage.settings.IdProperty) {
		return fmt.Errorf("it is not permitted to delete the id property ('%s')", info.storage.settings.IdProperty)
	}

	err := info.Set(ctx, lo.Map(properties, func(prop string, _ int) Entry { return Entry{prop, nil} })...)
	if err != nil {
		return fmt.Errorf("failed to set values to nil: %s", err)
	}

	return nil
}

func (info Info) ID() string {
	return info.id
}

func (info Info) Storage() *Storage {
	return info.storage
}

func (info Info) set() record.Set {
	return info.storage.records.Filter(record.E{info.storage.settings.IdProperty, info.id})
}

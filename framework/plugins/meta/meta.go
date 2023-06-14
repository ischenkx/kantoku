package meta

import (
	"context"
	"kantoku/common/codec"
	"kantoku/common/data/kv"
)

type RawMeta kv.Database[string, []byte]

type Meta struct {
	raw   RawMeta
	codec codec.Dynamic[[]byte]
}

type Value struct {
	name string
	meta Meta
}

func (meta Meta) Get(name string) Value {
	return Value{
		name: name,
		meta: meta,
	}
}

func (value Value) Set(ctx context.Context, val any) error {
	encoded, err := value.meta.codec.Encode(val)
	if err != nil {
		return err
	}

	return value.meta.raw.Set(ctx, value.name, encoded)
}

func (value Value) Load(ctx context.Context, to any) error {
	raw, err := value.LoadRaw(ctx)
	if err != nil {
		return err
	}
	return value.meta.codec.Decode(raw, to)
}

func (value Value) LoadRaw(ctx context.Context) ([]byte, error) {
	return value.meta.raw.Get(ctx, value.name)
}

func (value Value) Erase(ctx context.Context) error {
	return value.meta.raw.Del(ctx, value.name)
}

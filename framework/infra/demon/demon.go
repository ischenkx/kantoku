package demon

import (
	"context"
	"reflect"
)

type Demon struct {
	Type      string
	Name      string
	Parameter any
}

func (demon Demon) Eq(other Demon) bool {
	return reflect.DeepEqual(demon, other)
}

func (demon Demon) Demons(_ context.Context) []Demon {
	return []Demon{demon}
}

type Provider interface {
	Demons(ctx context.Context) []Demon
}

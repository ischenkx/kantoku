package infra

import (
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

func (demon Demon) Demons() []Demon {
	return []Demon{demon}
}

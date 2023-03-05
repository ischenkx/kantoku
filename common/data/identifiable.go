package data

import "kantoku/common/transformator"

type Identifiable interface {
	ID() string
}

type dynamicIdentifiable[T any] struct {
	item       *T
	identifier transformator.Transformator[*T, string]
}

func (d dynamicIdentifiable[T]) ID() string {
	return d.identifier(d.item)
}

func MakeIdentifiable[T any](item T, identifier transformator.Transformator[*T, string]) Identifiable {
	return dynamicIdentifiable[T]{
		item:       &item,
		identifier: identifier,
	}
}


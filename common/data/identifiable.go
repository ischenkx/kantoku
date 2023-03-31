package data

import "kantoku/common/transformer"

type Identifiable interface {
	ID() string
}

type dynamicIdentifiable[T any] struct {
	item       *T
	identifier transformer.Transformer[*T, string]
}

func (d dynamicIdentifiable[T]) ID() string {
	return d.identifier(d.item)
}

func MakeIdentifiable[T any](item T, identifier transformer.Transformer[*T, string]) Identifiable {
	return dynamicIdentifiable[T]{
		item:       &item,
		identifier: identifier,
	}
}

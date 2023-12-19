package bimap

import "context"

type Bimap[K, V any] interface {
	Set(context.Context, K, V) error
	DeleteByKey(context.Context, K) error
	DeleteByValue(context.Context, V) error
	ByValue(context.Context, V) (K, error)
	ByKey(context.Context, K) (V, error)
}

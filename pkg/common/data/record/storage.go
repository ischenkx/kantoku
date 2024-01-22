package record

import (
	"context"
)

type Storage[Item any] interface {
	Insert(context.Context, Item) error
	Set[Item]
}

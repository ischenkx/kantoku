package record

import (
	"context"
)

type Storage interface {
	Insert(context.Context, Record) error
	Set
}

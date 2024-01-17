package recutil

import (
	"context"
	"errors"
	"fmt"
	"github.com/ischenkx/kantoku/pkg/common/data/record"
)

func ListIter[Item any](ctx context.Context, iter record.Iter[Item]) ([]Item, error) {
	var result []Item

	for {
		item, err := iter.Next(ctx)
		if err != nil {
			if errors.Is(err, record.ErrIterEmpty) {
				break
			}

			return nil, fmt.Errorf("failed to iterate: %w", err)
		}
		result = append(result, item)
	}

	return result, nil
}

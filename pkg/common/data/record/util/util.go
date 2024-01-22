package recutil

import (
	"context"
	"errors"
	"fmt"
	"github.com/ischenkx/kantoku/pkg/common/data/record"
)

func List[Item any](ctx context.Context, iter record.Iter[Item]) ([]Item, error) {
	var result []Item

	defer iter.Close(ctx)
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

func Single[Item any](ctx context.Context, iter record.Iter[Item]) (Item, error) {
	defer iter.Close(ctx)

	result, err := iter.Next(ctx)
	if err != nil {
		return result, fmt.Errorf("failed to iterate: %w", err)
	}

	return result, nil
}

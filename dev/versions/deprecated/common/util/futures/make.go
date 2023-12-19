package futures

import (
	"context"
	future2 "kantoku/common/data/future"
)

func Make(ctx context.Context, manager *future2.Manager, ids ...*future2.ID) error {
	for _, id := range ids {
		res, err := manager.Make(ctx)
		if err != nil {
			return err

		}
		*id = res
	}
	return nil
}

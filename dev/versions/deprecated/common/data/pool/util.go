package pool

import (
	"context"
	"fmt"
	"kantoku/common/data/transactional"
)

func AutoCommit[Item any](
	ctx context.Context,
	reader Reader[Item],
	f func(ctx context.Context, item Item) error) error {

	ch, err := reader.Read(ctx)
	if err != nil {
		return fmt.Errorf("failed to open read channel: %s", err)
	}
loop:
	for {
		select {
		case <-ctx.Done():
			break loop
		case tx := <-ch:
			err = func(tx transactional.Object[Item]) error {
				defer tx.Rollback(ctx)

				item, err := tx.Get(ctx)
				if err != nil {
					return fmt.Errorf("failed to get item from transaction: %s", err)
				}

				if err := f(ctx, item); err != nil {
					return err
				}

				if err := tx.Commit(ctx); err != nil {
					return fmt.Errorf("failed to commit: %s", err)
				}

				return nil
			}(tx)

			if err != nil {
				return err
			}
		}
	}
	return nil
}

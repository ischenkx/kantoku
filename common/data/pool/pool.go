package pool

import (
	"context"
	"fmt"
	"kantoku/common/data/transaction"
)

//type Reader[Item any] interface {
//	Read(ctx context.Context) (<-chan Item, error)
//}
//
//type Writer[Item any] interface {
//	Write(ctx context.Context, item Item) error
//}
//
//type Pool[Item any] interface {
//	Reader[Item]
//	Writer[Item]
//}

type Reader[Item any] interface {
	Read(ctx context.Context) (<-chan transaction.Object[Item], error)
	//pool.Reader[Transaction[Item]]
}

// Writer probably should have NewTransaction method
type Writer[Item any] interface {
	// Write *must* write all items in a transaction!
	Write(ctx context.Context, items ...Item) error
	//pool.Writer[Item]
}

type Pool[Item any] interface {
	Reader[Item]
	Writer[Item]
}

func ReadAutoCommit[Item any](ctx context.Context, reader Reader[Item],
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
			err = func(tx transaction.Object[Item]) error {
				item, err := tx.Get(ctx)
				defer tx.Rollback(ctx)
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

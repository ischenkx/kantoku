package depot

import (
	"context"
	"fmt"
	"kantoku/framework/infra/demon"
	"kantoku/framework/utils/demons"
	"log"
)

type Plugin struct {
	depot *Depot
}

func NewPlugin(depot *Depot) Plugin {
	return Plugin{depot: depot}
}

func (plugin Plugin) Demons(ctx context.Context) []demon.Demon {
	return demons.Multi{
		demons.TryProvider(plugin.depot.Deps()),
		demons.TryProvider(plugin.depot.GroupTaskBimap()),
		demons.Functional("DEPOT_PROCESSOR", plugin.process),
	}.Demons(ctx)
}

func (plugin Plugin) process(ctx context.Context) error {
	ready, err := plugin.depot.Deps().Ready(ctx)
	if err != nil {
		return fmt.Errorf("failed to initialize the 'ready' channel: %s", err)
	}

loop:
	for {
		select {
		case <-ctx.Done():
			break loop
		case tx := <-ready:
			func() {
				log.Println(tx)
				defer tx.Rollback(ctx)
				id, err := tx.Get(ctx)
				if err != nil {
					log.Println("failed to get an id of a group:", err)
					return
				}
				taskID, err := plugin.depot.groupTaskBimap.ByKey(ctx, id)
				if err != nil {
					log.Println("failed to get a task assigned to the group:", err)
					return
				}
				if err := plugin.depot.inputs.Write(ctx, taskID); err != nil {
					log.Println("failed to schedule a task:", err)
					return
				}
				if err := tx.Commit(ctx); err != nil {
					log.Println("INCONSISTENCY! failed to commit reading group:", err)
				}
			}()
		}
	}

	return nil
}

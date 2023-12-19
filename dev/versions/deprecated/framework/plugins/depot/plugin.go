package depot

import (
	"context"
	"fmt"
	demons2 "kantoku/common/util/demons"
	"kantoku/framework/infra"
	"log"
)

type Plugin struct {
	depot *Depot
}

func NewPlugin(depot *Depot) Plugin {
	return Plugin{depot: depot}
}

func (plugin Plugin) Demons() []infra.Demon {
	return demons2.Multi{
		demons2.TryProvider(plugin.depot.Deps()),
		demons2.TryProvider(plugin.depot.GroupTaskBimap()),
		demons2.Functional("DEPOT_PROCESSOR", plugin.process),
	}.Demons()
}

func (plugin Plugin) process(ctx context.Context) error {
	log.Println("PROCESSING!")
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

package depot

import (
	"context"
	"errors"
	"fmt"
	"kantoku/common/data/bimap"
	"kantoku/common/data/transactional"
	"kantoku/framework/plugins/depot/deps"
	"kantoku/kernel"
	"kantoku/kernel/platform"
	"log"
)

type PluginData struct {
	Dependencies []string
}

type Depot struct {
	deps           deps.Deps
	groupTaskBimap bimap.Bimap[string, string]
	inputs         platform.Inputs
}

func New(deps deps.Deps, groupTaskBimap bimap.Bimap[string, string], inputs platform.Inputs) *Depot {
	return &Depot{
		deps:           deps,
		groupTaskBimap: groupTaskBimap,
		inputs:         inputs,
	}
}

func (depot *Depot) Deps() deps.Deps {
	return depot.deps
}

func (depot *Depot) GroupTaskBimap() bimap.Bimap[string, string] {
	return depot.groupTaskBimap
}

func (depot *Depot) Write(ctx context.Context, ids ...string) error {
	// actually it breaks interface
	if len(ids) != 1 {
		return fmt.Errorf("multiple ids are not supported")
	}
	taskId := ids[0]

	data := kernel.GetPluginData(ctx).GetWithDefault("dependencies", &PluginData{}).(*PluginData)

	// todo: somehow process error after intercepting
	group, err := depot.Deps().MakeGroup(ctx, func(ctx context.Context, id string) error {
		err := depot.groupTaskBimap.Save(ctx, id, taskId)
		if err != nil {
			return fmt.Errorf("failed to save the (group, task) pair in the bimap: %w", err)
		}
		return nil
	}, data.Dependencies...)
	if errors.Is(err, "failed to save the (group, task) pair in the bimap:")
	if err != nil {
		return fmt.Errorf("failed to make a dependency group: %s", err)
	}

	return nil
}

func (depot *Depot) Read(ctx context.Context) (<-chan transactional.Object[string], error) {
	return depot.inputs.Read(ctx)
}

func (depot *Depot) Process(ctx context.Context) error {
	ready, err := depot.Deps().Ready(ctx)
	if err != nil {
		return fmt.Errorf("failed to initialize the 'ready' channel: %s", err)
	}

loop:
	for {
		log.Println("reading from queue")
		select {
		case <-ctx.Done():
			break loop
		case tx := <-ready:
			func() {
				defer tx.Rollback(ctx)
				id, err := tx.Get(ctx)
				if err != nil {
					log.Println("failed to get an id of a group:", err)
					return
				}
				taskID, err := depot.groupTaskBimap.ByKey(ctx, id)
				if err != nil {
					log.Println("failed to get a task assigned to the group:", err)
					return
				}
				if err := depot.inputs.Write(ctx, taskID); err != nil {
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

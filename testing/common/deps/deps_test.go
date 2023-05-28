package deps

import (
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/samber/lo"
	"github.com/stretchr/testify/assert"
	"kantoku/common/deps"
	mempool "kantoku/impl/common/data/pool/mem"
	"kantoku/impl/common/deps/postgredeps"
	"math/rand"
	"testing"
	"time"
)

func newPostgresDeps(ctx context.Context) *postgredeps.Deps {
	client, err := pgxpool.New(ctx, "postgres://postgres:root@localhost:5432/")

	if err != nil {
		panic("failed to create postgres deps: " + err.Error())
	}

	if err := client.Ping(ctx); err != nil {
		panic("failed to make ping postgres: " + err.Error())
	}

	app := postgredeps.New(client, mempool.New[string]())

	app.Run(ctx)
	return app
}

func TestDeps(t *testing.T) {
	ctx := context.Background()
	implementations := map[string]deps.Deps{
		"postgres": newPostgresDeps(ctx),
	}

	for label, impl := range implementations {
		t.Run(label, func(t *testing.T) {
			dependencies := make([]string, rand.Intn(10))
			for i := 0; i < len(dependencies); i++ {
				dep, err := impl.Make(ctx)
				if err != nil {
					t.Fatal("failed to make a dependency:", err)
				}
				dependencies[i] = dep.ID
			}

			groupID, err := impl.MakeGroup(ctx, dependencies...)
			if err != nil {
				t.Fatal("failed to make a group:", err)
			}

			dep2resolution := lo.SliceToMap(dependencies, func(item string) (string, bool) { return item, false })

			for _, dep := range dependencies {
				if err := impl.Resolve(ctx, dep); err != nil {
					t.Fatalf("failed to resolve a dependency (group='%s', dep='%s'): %s", groupID, dep, err)
				}
				dep2resolution[dep] = true

				group, err := impl.Group(ctx, groupID)
				if err != nil {
					t.Fatalf("failed to get the group (%s): %s", groupID, err)
				}
				groupResolutions := lo.SliceToMap(group.Dependencies, func(item deps.Dependency) (string, bool) {
					return item.ID, item.Resolved
				})

				assert.Equal(t, group.ID, groupID, "group id is not equal to a retrieved group id")
				assert.Equal(t, dep2resolution, groupResolutions)
			}

			cancelableContext, cancel := context.WithCancel(ctx)
			defer cancel()

			ch, err := impl.Ready(cancelableContext)

			time.Sleep(time.Second * 3)

			select {
			case id := <-ch:
				if id != groupID {
					t.Fatalf("received a wrong group id: '%s' (expected '%s')", id, groupID)
				}
			default:
				t.Fatal("didn't receive a group from ready...")
			}
		})
	}
}

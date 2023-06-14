package deps

import (
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/samber/lo"
	"github.com/stretchr/testify/assert"
	deps2 "kantoku/framework/plugins/depot/deps"
	mempool "kantoku/impl/common/data/pool/mem"
	"kantoku/impl/plugins/deps/postgres"
	"testing"
	"time"
)

func newPostgresDeps(ctx context.Context) *postgredeps.Deps {
	client, err := pgxpool.New(ctx, "postgres://postgres:51413@localhost:5432/")

	if err != nil {
		panic("failed to create postgres deps: " + err.Error())
	}

	if err := client.Ping(ctx); err != nil {
		panic("failed to make ping postgres: " + err.Error())
	}

	app := postgredeps.New(client, mempool.New[string]())
	err = app.DropTables(ctx)
	if err != nil {
		panic("failed to init postgres tables: " + err.Error())
	}
	err = app.InitTables(ctx)
	if err != nil {
		panic("failed to init postgres tables: " + err.Error())
	}
	app.Run(ctx)
	return app
}

func TestDeps(t *testing.T) {
	ctx := context.Background()
	implementations := map[string]deps2.Deps{
		"postgres": newPostgresDeps(ctx),
	}

	for label, impl := range implementations {
		t.Run(label+"basic", func(t *testing.T) {
			dependencies, groupID := makeSimpleGroup(ctx, t, impl, 10)

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
				groupResolutions := lo.SliceToMap(group.Dependencies, func(item deps2.Dependency) (string, bool) {
					return item.ID, item.Resolved
				})

				assert.Equal(t, group.ID, groupID, "group id is not equal to a retrieved group id")
				assert.Equal(t, dep2resolution, groupResolutions)
			}

			cancelableContext, cancel := context.WithCancel(ctx)
			defer cancel()

			ch, err := impl.Ready(cancelableContext)
			if err != nil {
				t.Fatalf("failed to get a ready channel: %s", err)
			}

			checkReady(ch, func(id string) {
				if id != groupID {
					t.Fatalf("received a wrong group id: '%s' (expected '%s')", id, groupID)
				}
			}, func() {
				t.Fatal("didn't receive a group from ready...")
			})
		})

		t.Run(label+" double group counter decrement", func(t *testing.T) {
			dependencies, groupID := makeSimpleGroup(ctx, t, impl, 10)

			cancelableContext, cancel := context.WithCancel(ctx)
			defer cancel()
			ch, err := impl.Ready(cancelableContext)
			if err != nil {
				t.Fatalf("failed to get a ready channel: %s", err)
			}

			// resolve same dependency 10 times
			for i := 0; i < 10; i++ {
				err = impl.Resolve(ctx, dependencies[0])
				if err != nil {
					t.Fatalf("failed to resolve: %s", err)
				}
			}

			checkReady(ch, func(id string) {
				if id == groupID {
					t.Fatalf("group(%s) got to ready channel after resolving only one dependency(%s)",
						groupID, dependencies[0])
				} else {
					t.Fatalf("were expecting no ids, or %s in case of wrong behaviour, but received %s",
						groupID, id)
				}
			}, func() {})

			// resolve all dependencies	except last one
			for _, dep := range dependencies[0 : len(dependencies)-1] {
				for i := 0; i < 10; i++ {
					err = impl.Resolve(ctx, dep)
					if err != nil {
						t.Fatalf("failed to resolve: %s", err)
					}
				}
			}
		})
	}
}

// create channel here?
func checkReady(ch <-chan string, receive func(id string), nothing func()) {
	time.Sleep(time.Second * 3)
	select {
	case id := <-ch:
		receive(id)
	default:
		nothing()
	}
}

func makeSimpleGroup(ctx context.Context, t *testing.T, impl deps2.Deps, size int) ([]string, string) {
	dependencies := make([]string, size)
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
	return dependencies, groupID
}

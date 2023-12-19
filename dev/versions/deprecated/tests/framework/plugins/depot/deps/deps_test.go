package deps

import (
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/samber/lo"
	"github.com/stretchr/testify/assert"
	deps2 "kantoku/common/data/deps"
	"kantoku/common/data/transactional"
	mempool "kantoku/impl/common/data/pool/mem"
	"kantoku/impl/deps/postgres/batched"
	"kantoku/impl/deps/postgres/instant"
	"log"
	"strings"
	"testing"
	"time"
)

func newBatchedDeps(ctx context.Context) deps2.Deps {
	client, err := pgxpool.New(ctx, "postgres://postgres:51413@localhost:5432/")

	if err != nil {
		panic("failed to create postgres deps: " + err.Error())
	}

	if err := client.Ping(ctx); err != nil {
		panic("failed to make ping postgres: " + err.Error())
	}

	app := batched.New(client, mempool.New[string](mempool.DefaultConfig))
	err = app.DropTables(ctx)
	if err != nil {
		log.Println("Warning: failed to drop postgres tables: ", err.Error())
	}
	err = app.InitTables(ctx)
	if err != nil {
		panic("failed to init postgres tables: " + err.Error())
	}
	app.Run(ctx)
	return app
}

func newInstantDeps(ctx context.Context) deps2.Deps {
	client, err := pgxpool.New(ctx, "postgres://postgres:51413@localhost:5432/")

	if err != nil {
		panic("failed to create postgres deps: " + err.Error())
	}

	if err := client.Ping(ctx); err != nil {
		panic("failed to make ping postgres: " + err.Error())
	}

	app := instant.New(client, mempool.New[string](mempool.DefaultConfig))
	err = app.DropTables(ctx)
	if err != nil {
		log.Println("Warning: failed to drop postgres tables: ", err.Error())
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
	implementations := map[string]func(context.Context) deps2.Deps{
		"postgres-batched": newBatchedDeps,
		"postgres-instant": newInstantDeps,
	}

	for label, newImpl := range implementations {
		t.Run(label, func(t *testing.T) {
			t.Run("basic", func(t *testing.T) {
				ctx, cancel := context.WithCancel(ctx)
				defer cancel()

				impl := newImpl(ctx)
				dependencies, groupID := makeSimpleGroup(ctx, impl, t, 10)

				dep2resolution := lo.SliceToMap(dependencies, func(item string) (string, bool) { return item, false })

				for _, dep := range dependencies {
					resolveDep(ctx, impl, t, dep)
					dep2resolution[dep] = true

					group := getGroup(ctx, impl, t, groupID)
					groupResolutions := lo.SliceToMap(group.Dependencies, func(item deps2.Dependency) (string, bool) {
						return item.ID, item.Resolved
					})

					assert.Equal(t, group.ID, groupID, "group id is not equal to a retrieved group id")
					assert.Equal(t, dep2resolution, groupResolutions)
				}

				ch := getReadyCh(ctx, impl, t)
				checkReady(ctx, t, ch, func(id string) {
					if id != groupID {
						t.Fatalf("received a wrong group id: '%s' (expected '%s')", id, groupID)
					}
				}, func() {
					t.Fatalf("didn't receive a group from ready, expected '%s'", groupID)
				})
			})

			t.Run("double group counter increment", func(t *testing.T) {
				ctx, cancel := context.WithCancel(ctx)
				defer cancel()

				impl := newImpl(ctx)
				dependencies, groupID := makeSimpleGroup(ctx, impl, t, 10)

				ch := getReadyCh(ctx, impl, t)
				// resolve same dependency 10 times
				for i := 0; i < 10; i++ {
					resolveDep(ctx, impl, t, dependencies[0])
				}

				checkReady(ctx, t, ch, func(id string) {
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
						resolveDep(ctx, impl, t, dep)
					}
				}

				checkReady(ctx, t, ch, func(id string) {
					if id == groupID {
						t.Fatalf("group(%s) got to ready channel after resolving not all of dependencies (%s - not resolved)",
							groupID, dependencies[len(dependencies)-1])
					} else {
						t.Fatalf("were expecting no ids, or %s in case of wrong behaviour, but received %s",
							groupID, id)
					}
				}, func() {})
			})

			t.Run("group from resolved deps", func(t *testing.T) {
				ctx, cancel := context.WithCancel(ctx)
				defer cancel()

				impl := newImpl(ctx)
				ready := getReadyCh(ctx, impl, t)

				for i := 1; i < 10; i++ {
					dependencies, skip := makeSimpleGroup(ctx, impl, t, i)
					for _, dep := range dependencies {
						resolveDep(ctx, impl, t, dep)
						depData := getDep(ctx, impl, t, dep)
						assert.True(t, depData.Resolved, "dependency hasn't resolved (%s)", dep)
					}
					// skip normal group
					checkReady(ctx, t, ready, func(id string) {
						assert.Equal(t, skip, id, "received unexpected group")
					}, func() {
						t.Fatalf("haven't received group (%s), but all it's dependencies(%s) are resolved",
							skip, strings.Join(dependencies, ", "))
					})

					group := makeGroup(ctx, impl, t, dependencies...)
					checkReady(ctx, t, ready, func(id string) {
						assert.Equal(t, group, id, "received unexpected group")
					}, func() {
						t.Fatalf("haven't received group (%s), but all it's dependencies(%s) are resolved",
							group, strings.Join(dependencies, ", "))
					})
				}
			})
		})
	}
}

// create channel here?
func checkReady(ctx context.Context, t *testing.T, ch <-chan transactional.Object[string], receive func(id string), nothing func()) {
	select {
	case tx := <-ch:
		id, err := tx.Get(ctx)
		if err != nil {
			t.Fatalf("failed to get value of transaction: %s", err)
		}

		receive(id)

		if err := tx.Commit(ctx); err != nil {
			t.Fatalf("failed to commit a transaction: %s", err)
		}
	case <-time.After(7 * time.Second):
		nothing()
	}
}

func makeSimpleGroup(ctx context.Context, impl deps2.Deps, t *testing.T, size int) ([]string, string) {
	dependencies := make([]string, size)
	for i := 0; i < len(dependencies); i++ {
		dep := makeDep(ctx, impl, t)
		dependencies[i] = dep.ID
	}
	groupID := makeGroup(ctx, impl, t, dependencies...)
	return dependencies, groupID
}

func getGroup(ctx context.Context, impl deps2.Deps, t *testing.T, id string) deps2.Group {
	group, err := impl.Group(ctx, id)
	if err != nil {
		t.Fatalf("failed to get a group(%s): %s", id, err)
	}
	return group
}

func getDep(ctx context.Context, impl deps2.Deps, t *testing.T, id string) deps2.Dependency {
	dep, err := impl.Dependency(ctx, id)
	if err != nil {
		t.Fatalf("failed to get a dependecy(%s): %s", id, err)
	}
	return dep
}

func makeDep(ctx context.Context, impl deps2.Deps, t *testing.T) deps2.Dependency {
	dep, err := impl.NewDependency(ctx)
	if err != nil {
		t.Fatal("failed to make a dependency:", err)
	}
	return dep
}

func makeGroup(ctx context.Context, impl deps2.Deps, t *testing.T, ids ...string) string {
	group, err := impl.NewGroup(ctx)
	if err != nil {
		t.Fatalf("failed to create group from dependecies(%s):\n%s", strings.Join(ids, ", "), err)
	}
	if err := impl.InitGroup(ctx, group, ids...); err != nil {
		t.Fatalf("failed to create group from dependecies(%s):\n%s", strings.Join(ids, ", "), err)
	}
	return group
}

func resolveDep(ctx context.Context, impl deps2.Deps, t *testing.T, id string) {
	err := impl.Resolve(ctx, id)
	if err != nil {
		t.Fatalf("failed to resolve dependecy(%s): %s", id, err)
	}
}

func getReadyCh(ctx context.Context, impl deps2.Deps, t *testing.T) <-chan transactional.Object[string] {
	ch, err := impl.Ready(ctx)
	if err != nil {
		t.Fatalf("failed to get a ready channel: %s", err)
	}
	return ch
}

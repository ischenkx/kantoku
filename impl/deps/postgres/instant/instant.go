package instant

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"kantoku/common/data/pool"
	"kantoku/common/data/transactional"
	"kantoku/framework/plugins/depot/deps"
	"kantoku/impl/deps/postgres"
)

var _ postgres.Deps = &Deps{}

type Deps struct {
	resolvedGroups pool.Pool[string]
	client         *pgxpool.Pool
}

// New creates an instance of the dependency manager with *instant* resolving.
// client is a connection to a postgres database (all dependencies and groups are stored there)
// queue is for storing resolved groups waiting to be processed
func New(client *pgxpool.Pool, groupsQueue pool.Pool[string]) *Deps {
	return &Deps{
		resolvedGroups: groupsQueue,
		client:         client,
	}
}

func (d *Deps) InitTables(ctx context.Context) error {
	sql := `
		CREATE TABLE InstantDependencies (
			id varchar(255),
			resolved bool,
			PRIMARY KEY (id)
		);
		
		CREATE TABLE InstantGroups (
			id varchar(255),
			pending int,
			PRIMARY KEY (id)
		);
		
		CREATE TABLE InstantGroupDependencies (
			dependency_id varchar(255),
			group_id varchar(255),
			PRIMARY KEY (dependency_id, group_id),
			CONSTRAINT fk_group
			  FOREIGN KEY(group_id) 
			  REFERENCES InstantGroups(id)
		      ON DELETE CASCADE
		);
	`
	_, err := d.client.Exec(ctx, sql)
	return err
}

func (d *Deps) Dependency(ctx context.Context, id string) (deps.Dependency, error) {
	sql := `SELECT id, resolved FROM InstantDependencies WHERE id = $1`
	row := d.client.QueryRow(ctx, sql, id)

	var model deps.Dependency
	if err := row.Scan(&model.ID, &model.Resolved); err != nil {
		return deps.Dependency{}, err
	}
	return model, nil
}

func (d *Deps) DropTables(ctx context.Context) error {
	sql := `DROP TABLE InstantDependencies, InstantGroups, InstantGroupDependencies;`
	_, err := d.client.Exec(ctx, sql)
	return err
}

func (d *Deps) Group(ctx context.Context, group string) (deps.Group, error) {
	var result deps.Group
	result.ID = group

	sql := `SELECT dependency_id, COALESCE(resolved, false) FROM InstantGroupDependencies gd 
				LEFT JOIN InstantDependencies d ON d.id = gd.dependency_id
				WHERE gd.group_id = $1`
	records, err := d.client.Query(ctx, sql, group)
	if err != nil {
		return result, err
	}
	defer records.Close()

	for records.Next() {
		var dep deps.Dependency
		if err := records.Scan(&dep.ID, &dep.Resolved); err != nil {
			return result, err
		}
		result.Dependencies = append(result.Dependencies, dep)
	}

	return result, nil
}

// Make generates a new dependency (but it does not store the information about it in the database
//
// Generation algorithm: UUID
func (d *Deps) Make(_ context.Context) (deps.Dependency, error) {
	id := uuid.New().String()

	return deps.Dependency{
		ID:       id,
		Resolved: false,
	}, nil
}

// MakeGroup creates a new dependency group
//
// NOTE: the new group's id is generated via a UUID algorithm
func (d *Deps) MakeGroup(ctx context.Context, intercept func(context.Context, string) error,
	ids ...string) (string, error) {

	id := uuid.New().String()

	if err := intercept(ctx, id); err != nil {
		return "", fmt.Errorf("failed to call intercept: %w", err)
	}

	tx, err := d.client.Begin(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to begin a postgres transaction: %s", tx)
	}
	defer tx.Rollback(ctx)

	sql := `
			INSERT INTO InstantGroups (id, pending)
			VALUES ($1, $2 - (SELECT COUNT(*) FROM InstantDependencies as d WHERE d.id = ANY ($3)))
			RETURNING pending
			`
	result := tx.QueryRow(ctx, sql, id, len(ids), ids)
	var pending int
	if err := result.Scan(&pending); err != nil {
		return "", fmt.Errorf("failed to scan number of pending dependencies for a group(%s): %s", id, err)
	}
	// if pending = 0, we need to add group to pool. It is done in the end of the method, after committing transaction

	for _, dep := range ids {
		sql := `INSERT INTO InstantGroupDependencies (dependency_id, group_id) VALUES ($1, $2)`
		_, err := tx.Exec(ctx, sql, dep, id)
		if err != nil {
			return "", err
		}
	}

	if err := tx.Commit(ctx); err != nil {
		return "", fmt.Errorf("failed to commit the transaction: %s", err)
	}
	// it is possible that adding to pool but if it is so, it's id won't be used, and it shouldn't break anything
	if pending == 0 {
		// all deps have already been resolved
		if err := d.resolvedGroups.Write(ctx, id); err != nil {
			return "", fmt.Errorf("failed to add resolved group to pool: %s", err)
		}
	}

	return id, nil
}

func (d *Deps) Resolve(ctx context.Context, dep string) error {
	//if err := d.resolvedDeps.Write(ctx, dep); err != nil {
	//	return fmt.Errorf("failed to write dep to queue: %s", err)
	//}
	//return nil

	tx, err := d.client.Begin(ctx)
	if err != nil {
		return fmt.Errorf("failed to begin a postgres transaction: %s", tx)
	}
	defer tx.Rollback(ctx)
	// let's goooooo
	// 1) update dep and check if it has been resolved
	sql := `
		WITH previous AS (
			SELECT resolved
			FROM InstantDependencies
			WHERE id = $1
			LIMIT 1
		)
		INSERT INTO InstantDependencies (id, resolved) 
		VALUES ($1, $2)
		ON CONFLICT (id) DO UPDATE 
		SET resolved = $2
		RETURNING COALESCE((SELECT resolved FROM previous), false) AS previous_resolved;`
	before := tx.QueryRow(ctx, sql, dep, true)
	var alreadyResolved bool
	if err := before.Scan(&alreadyResolved); err != nil {
		return err
	}
	if alreadyResolved {
		// if it resolved we are done
		return nil
	}

	// dep hasn't been resolved, so we need to update groups
	// we also return groups which have no unresolved deps
	sql = `
			UPDATE InstantGroups g
			SET pending = pending - 1
			FROM InstantGroupDependencies gd
			WHERE gd.dependency_id = $1 AND gd.group_id = g.id
			RETURNING (SELECT g.id WHERE g.pending = 0)` // returns new value
	resolved, err := tx.Query(ctx, sql, dep)
	if err != nil {
		return fmt.Errorf("failed to decrement group counters: %s", err)
	}
	defer resolved.Close()

	var groups []string
	for resolved.Next() {
		// I do not understand wny it returns nulls when result should be empty
		var id *string
		if err := resolved.Scan(&id); err != nil {
			return fmt.Errorf("failed to scan the id of a group (please report if it ever happens): %s", err)
		}
		if id != nil {
			groups = append(groups, *id)
		}
	}
	// now we just add groups to their queue
	if err := d.resolvedGroups.Write(ctx, groups...); err != nil {
		return fmt.Errorf("failed to add groups to queue: %s", err)
	}

	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("failed to commit the transaction: %s", err)
	}
	return nil
}

func (d *Deps) Ready(ctx context.Context) (<-chan transactional.Object[string], error) {
	return d.resolvedGroups.Read(ctx)
}

func (d *Deps) Run(_ context.Context) {}

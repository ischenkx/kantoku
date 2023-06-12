package instant

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"kantoku/common/data/pool"
	"kantoku/common/data/transaction"
	"kantoku/impl/deps/postgres"
	"kantoku/unused/backend/framework/depot/deps"
	"log"
)

var _ postgres.Deps = &Deps{}

type Deps struct {
	resolvedGroups pool.Pool[string]
	resolvedDeps   pool.Pool[string]
	client         *pgxpool.Pool
}

// New creates an instance of the dependency manager with *instant* resolving.
// client is a connection to a postgres database (all dependencies and groups are stored there)
// queues are for storing resolved groups and deps, waiting to be processed
func New(client *pgxpool.Pool, groupsQueue, depsQueue pool.Pool[string]) *Deps {
	return &Deps{
		resolvedGroups: groupsQueue,
		resolvedDeps:   depsQueue,
		client:         client,
	}
}

func (d *Deps) InitTables(ctx context.Context) error {
	sql := `
		CREATE TABLE Dependencies (
			id varchar(255),
			resolved bool,
			PRIMARY KEY (id)
		);
		
		CREATE TABLE Groups (
			id varchar(255),
			pending int,
			PRIMARY KEY (id)
		);
		
		CREATE TABLE GroupDependencies (
			dependency_id varchar(255),
			group_id varchar(255),
			PRIMARY KEY (dependency_id, group_id),
			CONSTRAINT fk_group
			  FOREIGN KEY(group_id) 
			  REFERENCES Groups(id)
		      ON DELETE CASCADE
		);
	`
	_, err := d.client.Exec(ctx, sql)
	return err
}

func (d *Deps) Dependency(ctx context.Context, id string) (deps.Dependency, error) {
	sql := `SELECT id, resolved FROM Dependencies WHERE id = $1`
	row := d.client.QueryRow(ctx, sql, id)

	var model deps.Dependency
	if err := row.Scan(&model.ID, &model.Resolved); err != nil {
		return deps.Dependency{}, err
	}
	return model, nil
}

func (d *Deps) DropTables(ctx context.Context) error {
	sql := `DROP TABLE Dependencies, Groups, GroupDependencies;`
	_, err := d.client.Exec(ctx, sql)
	return err
}

func (d *Deps) Group(ctx context.Context, group string) (deps.Group, error) {
	var result deps.Group
	result.ID = group

	sql := `SELECT dependency_id, COALESCE(resolved, false) FROM GroupDependencies gd 
				LEFT JOIN Dependencies d ON d.id = gd.dependency_id
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
func (d *Deps) Make(ctx context.Context) (deps.Dependency, error) {
	id := uuid.New().String()

	return deps.Dependency{
		ID:       id,
		Resolved: false,
	}, nil
}

// MakeGroup creates a new dependency group
//
// NOTE: the new group's id is generated via a UUID algorithm
func (d *Deps) MakeGroup(ctx context.Context, ids ...string) (string, error) {
	id := uuid.New().String()

	tx, err := d.client.Begin(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to begin a postgres transaction: %s", tx)
	}
	defer tx.Rollback(ctx)

	sql := `INSERT INTO Groups (id, pending)
			VALUES ($1,
			        $2 - (SELECT COUNT(*) FROM dependencies d WHERE d.resolved)
			)`
	_, err = tx.Exec(ctx, sql, id, len(ids))
	if err != nil {
		return "", err
	}

	for _, dep := range ids {
		sql := `INSERT INTO GroupDependencies (dependency_id, group_id) VALUES ($1, $2)`
		_, err := tx.Exec(ctx, sql, dep, id)
		if err != nil {
			return "", err
		}
	}

	if err := tx.Commit(ctx); err != nil {
		return "", fmt.Errorf("failed to commit the transaction: %s", err)
	}

	return id, nil
}

func (d *Deps) Resolve(ctx context.Context, dep string) error {
	if err := d.resolvedDeps.Write(ctx, dep); err != nil {
		return fmt.Errorf("failed to write dep to queue: %s", err)
	}
	return nil
}

func (d *Deps) Ready(ctx context.Context) (<-chan transaction.Object[string], error) {
	return d.resolvedGroups.Read(ctx)
}

func (d *Deps) Run(ctx context.Context) {
	go d.runDepsProcessing(ctx)
}

func (d *Deps) runDepsProcessing(ctx context.Context) {
	err := pool.ReadAutoCommit[string](ctx, d.resolvedDeps, d.processDep)
	if err != nil {
		log.Println("failed to open a resolved deps channel:", err)
	}
}

func (d *Deps) processDep(ctx context.Context, dep string) error {
	sql := `
		INSERT INTO ProcessingDependencies (id)
		VALUES ($1)`

	if _, err := d.client.Exec(ctx, sql, dep); err != nil {
		return fmt.Errorf("failed to register dep processing: %s", err)
	}

	tx, err := d.client.Begin(ctx)
	if err != nil {
		return fmt.Errorf("failed to begin a postgres transaction: %s", tx)
	}
	defer tx.Rollback(ctx)
	// let's goooooo
	// 1) update dep and check if it has been resolved
	sql = `
		WITH previous AS (
			SELECT resolved
			FROM Dependencies
			WHERE id = $1
			LIMIT 1
		)
		INSERT INTO Dependencies (id, resolved) 
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
			UPDATE Groups g
			SET pending = pending - 1
			FROM GroupDependencies gd
			WHERE gd.dependency_id = $1 AND gd.group_id = g.id
			RETURNING (SELECT g.id WHERE g.pending = 0)` // returns new value
	resolved, err := tx.Query(ctx, sql, dep)
	var groups []string
	for resolved.Next() {
		var id string
		if err := resolved.Scan(&id); err != nil {
			return fmt.Errorf("failed to scan the id of a group (please report if it ever happens): %s", err)
		}
		groups = append(groups, id)
	}
	// now we just add groups to their queue
	if err := d.resolvedGroups.Write(ctx, groups...); err != nil {
		return fmt.Errorf("failed to add groups to queue")
	}

	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("failed to commit the transaction: %s", err)
	}
	return nil
}

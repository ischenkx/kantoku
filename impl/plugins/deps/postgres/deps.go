package postgredeps

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/samber/lo"
	"kantoku/common/data/pool"
	deps2 "kantoku/framework/plugins/depot/deps"
	"log"
	"strings"
	"time"
)

type Deps struct {
	q      pool.Pool[string]
	client *pgxpool.Pool
}

// New creates an instance of the dependency manager.
//
// client - a connection to a postgres database (all dependencies and groups are stored there)
//
// queue - a queue for "ready" groups
func New(client *pgxpool.Pool, queue pool.Pool[string]) *Deps {
	return &Deps{
		q:      queue,
		client: client,
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
			status varchar(16),
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

func (d *Deps) DropTables(ctx context.Context) error {
	sql := `DROP TABLE Dependencies, Groups, GroupDependencies;`
	_, err := d.client.Exec(ctx, sql)
	return err
}

func (d *Deps) Dependency(ctx context.Context, id string) (deps2.Dependency, error) {
	sql := `SELECT id, resolved FROM Dependencies WHERE id = $1`
	row := d.client.QueryRow(ctx, sql, id)

	var model deps2.Dependency
	if err := row.Scan(&model.ID, &model.Resolved); err != nil {
		return deps2.Dependency{}, err
	}
	return model, nil
}

func (d *Deps) Group(ctx context.Context, group string) (deps2.Group, error) {
	var result deps2.Group
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
		var dep deps2.Dependency
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
func (d *Deps) Make(ctx context.Context) (deps2.Dependency, error) {
	id := uuid.New().String()

	return deps2.Dependency{
		ID:       id,
		Resolved: false,
	}, nil
}

// MakeGroup creates a new dependency group
//
// NOTE: the new group's id is generated via a UUID algorithm
func (d *Deps) MakeGroup(ctx context.Context, ids ...string) (string, error) {
	id := uuid.New().String()
	status := InitializingStatus

	tx, err := d.client.Begin(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to begin a postgres transaction: %s", tx)
	}
	defer tx.Rollback(ctx)

	sql := `INSERT INTO Groups (id, pending, status) VALUES ($1, $2, $3)`
	_, err = tx.Exec(ctx, sql, id, -1, status)
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
	tx, err := d.client.Begin(ctx)
	if err != nil {
		return fmt.Errorf("failed to begin a postgres transaction: %s", tx)
	}
	defer tx.Rollback(ctx)

	sql := `
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
		return nil
	}

	sql = `
			UPDATE Groups g
			SET pending = pending - 1
			FROM GroupDependencies gd
			WHERE gd.dependency_id = $1 AND gd.group_id = g.id`
	if _, err := tx.Exec(ctx, sql, dep); err != nil {
		return err
	}

	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("failed to commit the transaction: %s", err)
	}
	return nil
}

func (d *Deps) Ready(ctx context.Context) (<-chan string, error) {
	return d.q.Read(ctx)
}

func (d *Deps) Run(ctx context.Context) {
	go d.runGroupsScheduler(ctx, time.Second, 1024)
	go d.runGroupInitializer(ctx, time.Second)
}

func (d *Deps) runGroupInitializer(ctx context.Context, interval time.Duration) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

loop:
	for {
		select {
		case <-ctx.Done():
			break loop
		case <-ticker.C:
			if err := d.initializeGroups(ctx); err != nil {
				log.Println("failed to initialize groups:", err)
			}
		}
	}
}

func (d *Deps) runGroupsScheduler(ctx context.Context, interval time.Duration, batchSize int) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

loop:
	for {
		select {
		case <-ctx.Done():
			break loop
		case <-ticker.C:
			if err := d.scheduleGroups(ctx, batchSize); err != nil {
				log.Println("failed to schedule groups:", err)
			}
		}
	}
}

func (d *Deps) initializeGroups(ctx context.Context) error {
	sql := `
		WITH uninitialized_groups AS (
		   SELECT id
		   FROM   Groups
		   WHERE  status = $1
		   ) 
		UPDATE Groups g
		SET status = $2, pending = (
			SELECT COUNT(*)
				FROM GroupDependencies gd
				LEFT OUTER JOIN Dependencies d ON d.id = gd.dependency_id
				WHERE gd.group_id = g.id
					AND CASE WHEN d.resolved is NULL THEN true ELSE NOT d.resolved END
		)
		FROM   uninitialized_groups
		WHERE  g.id = uninitialized_groups.id
		RETURNING g.id
	`

	_, err := d.client.Exec(ctx, sql, InitializingStatus, WaitingStatus)
	return err
}

func (d *Deps) scheduleGroups(ctx context.Context, batchSize int) error {
	sql := `
		WITH groups_subset AS (
		   SELECT id
		   FROM   Groups
		   WHERE  status = $1 AND pending = 0
		   LIMIT  $2
		   )
		UPDATE Groups g
		SET    status = $3 
		FROM   groups_subset
		WHERE  g.id = groups_subset.id
		RETURNING g.id
	`

	records, err := d.client.Query(ctx, sql, WaitingStatus, batchSize, SchedulingStatus)
	if err != nil {
		return err
	}

	var groups []string
	for records.Next() {
		var id string
		if err := records.Scan(&id); err != nil {
			log.Println("failed to scan the id of a group:", err)
			continue
		}

		groups = append(groups, id)
	}
	records.Close()

	var failed, succeeded []string
	for _, group := range groups {
		if err := d.q.Write(ctx, group); err != nil {
			failed = append(failed, group)
		} else {
			succeeded = append(succeeded, group)
		}
	}

	sql = `
		UPDATE Groups g
		SET    status = $1
		FROM   (values ($2)) as s(id)
		WHERE  g.id = s.id
	`

	formattedFailed := formatPostgresValues(failed...)
	if _, err := d.client.Exec(ctx, sql, WaitingStatus, formattedFailed); err != nil {
		log.Println("failed to update the groups that failed to be scheduled:", err)
	}

	formattedSucceeded := formatPostgresValues(succeeded...)
	if _, err := d.client.Exec(ctx, sql, ScheduledStatus, formattedSucceeded); err != nil {
		log.Println("failed to update the groups that were scheduled:", err)
	}
	return nil
}

func formatPostgresValues(ids ...string) string {
	values := strings.Join(
		lo.Map(ids, func(item string, _ int) string {
			return fmt.Sprintf("'%s'", item)
		}),
		", ")
	values = fmt.Sprintf("(%s)", values)
	return values
}

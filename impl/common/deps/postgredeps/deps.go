package postgredeps

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/samber/lo"
	"kantoku/common/deps"
	"kantoku/common/pool"
	"log"
	"strings"
	"time"
)

type Deps struct {
	q      pool.Pool[string]
	client *pgxpool.Conn
}

// MAKE DEPENDENCIES MONOTONOUS AND UNIQUE

func New(client *pgxpool.Conn, queue pool.Pool[string]) *Deps {
	return &Deps{
		q:      queue,
		client: client,
	}
}

func (d *Deps) InitTables(ctx context.Context) error {
	sql := `
		CREATE TABLE Dependencies (
			id varchar(255),
			last_resolution int,
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
			resolution int,
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

func (d *Deps) Dependency(ctx context.Context, dep string) (deps.Dependency, error) {
	sql := `SELECT id, last_resolution FROM Dependencies WHERE id = $1`
	row := d.client.QueryRow(ctx, sql, dep)

	var model deps.Dependency
	if err := row.Scan(&model.ID, &model.LastResolution); err != nil {
		return deps.Dependency{}, err
	}
	return model, nil
}

func (d *Deps) Group(ctx context.Context, group string) (deps.Group, error) {
	var result deps.Group
	result.ID = group
	result.Dependencies = map[string]bool{}

	sql := `SELECT dependency_id, resolution FROM GroupDependencies WHERE group_id = $1`
	records, err := d.client.Query(ctx, sql, group)
	if err != nil {
		return result, err
	}
	defer records.Close()

	for records.Next() {
		var dep string
		var resolution int
		if err := records.Scan(&dep, &resolution); err != nil {
			return result, err
		}
		result.Dependencies[dep] = resolution > 0
	}

	return result, nil
}

func (d *Deps) Make(ctx context.Context, ids ...string) (deps.Group, error) {
	id := uuid.New().String()
	status := WaitingStatus

	sql := `INSERT INTO Groups (id, pending, status) VALUES ($1, $2, $3)`
	_, err := d.client.Exec(ctx, sql, id, len(ids), status)
	if err != nil {
		return deps.Group{}, err
	}

	for _, dep := range ids {
		sql := `INSERT INTO GroupDependencies (dependency_id, group_id, resolution) VALUES ($1, $2, $3)`
		_, err := d.client.Exec(ctx, sql, dep, dep, id, -1)
		if err != nil {
			return deps.Group{}, err
		}
	}

	return deps.Group{
		ID: id,
		Dependencies: lo.Associate(ids, func(dep string) (string, bool) {
			return dep, false
		}),
	}, nil
}

func (d *Deps) Resolve(ctx context.Context, dep string) error {
	// Updating the last_resolution field in "Dependencies"
	timestamp := time.Now().UnixNano()
	sql := `
		INSERT INTO Dependencies (id, last_resolution) 
		VALUES ($1, $2)
		ON CONFLICT (id) DO UPDATE 
		SET last_resolution = $2`
	if _, err := d.client.Exec(ctx, sql, dep, timestamp); err != nil {
		return err
	}

	// Updating the group dependencies statuses
	sql = `UPDATE GroupDependencies SET resolution_ts = $1 WHERE dependency_id = $2`
	if _, err := d.client.Exec(ctx, sql, dep); err != nil {
		return err
	}

	// Updating the counters
	sql = `
			UPDATE Groups g
			SET pending_deps = pending_deps - 1
			FROM GroupDependencies gd
			WHERE gd.dependency_id = $1 AND gd.group_id = g.id AND gd.resolution = $2`
	if _, err := d.client.Exec(ctx, sql, dep, timestamp); err != nil {
		return err
	}

	return nil
}

func (d *Deps) Ready(ctx context.Context) (<-chan string, error) {
	return d.q.Read(ctx)
}

func (d *Deps) Run(ctx context.Context) {
	d.runGroupsScheduler(ctx, time.Second, 1024)
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
	defer records.Close()

	var groups []string

	for records.Next() {
		var id string
		if err := records.Scan(&id); err != nil {
			log.Println("failed to scan the id of a group:", err)
			continue
		}

		groups = append(groups, id)
	}

	var failed, succeeded []string

	for _, group := range groups {
		log.Println("scheduling:", group)
		if err := d.q.Write(ctx, group); err != nil {
			log.Println("failed:", group)
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

	formattedFailed := fmt.Sprintf("(%s)",
		strings.Join(
			lo.Map(failed, func(item string, index int) string {
				return fmt.Sprintf("'%s'", item)
			}),
			", "))

	formattedSucceeded := fmt.Sprintf("(%s)",
		strings.Join(
			lo.Map(succeeded, func(item string, index int) string {
				return fmt.Sprintf("'%s'", item)
			}),
			", "))
	if _, err := d.client.Exec(ctx, sql, WaitingStatus, formattedFailed); err != nil {
		log.Println("failed to update the groups that failed to be scheduled:", err)
	}

	if _, err := d.client.Exec(ctx, sql, ScheduledStatus, formattedSucceeded); err != nil {
		log.Println("failed to update the groups that were scheduled:", err)
	}
	return nil
}

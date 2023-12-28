package batched

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"github.com/ischenkx/kantoku/pkg/common/data/deps"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/samber/lo"
	"log/slog"
	"strings"
	"time"
)

var _ deps.Manager = (*Manager)(nil)

type Manager struct {
	client *pgxpool.Pool
	config Config
}

// New creates an instance of the dependency manager with *batched* resolving.
// client - a connection to a postgres database (all dependencies and groups are stored there)
// queue - a queue for "ready" groups
func New(client *pgxpool.Pool, config Config) *Manager {
	return &Manager{
		client: client,
		config: config,
	}
}

func (manager *Manager) LoadDependencies(ctx context.Context, ids ...string) ([]deps.Dependency, error) {
	sql := `select id, status from dependencies where id in ($1)`

	rows, err := manager.client.Query(ctx, sql, ids)
	if err != nil {
		return nil, fmt.Errorf("failed to query: %w", err)
	}

	result := make([]deps.Dependency, 0, len(ids))

	for rows.Next() {
		var dep deps.Dependency
		if err := rows.Scan(&dep.ID, &dep.Status); err != nil {
			return nil, fmt.Errorf("failed to scan: %w", err)
		}
		result = append(result, dep)
	}

	return result, nil
}

func (manager *Manager) LoadGroups(ctx context.Context, ids ...string) ([]deps.Group, error) {
	sql := `
			select group_id, dependency_id, d.status
			from group_dependencies gd
			join dependencies d on d.id = gd.dependency_id
			where group_id in ($1)
	`

	rows, err := manager.client.Query(ctx, sql, ids)
	if err != nil {
		return nil, fmt.Errorf("failed to query: %w", err)
	}

	result := make(map[string]deps.Group, len(ids))

	for rows.Next() {
		var groupId, dependencyId, status string

		if err := rows.Scan(&groupId, &dependencyId, &status); err != nil {
			return nil, fmt.Errorf("failed to scane: %w", err)
		}
		if _, ok := result[groupId]; !ok {
			result[groupId] = deps.Group{ID: groupId}
		}
		group := result[groupId]

		group.Dependencies = append(group.Dependencies, deps.Dependency{
			ID:     dependencyId,
			Status: deps.Status(status),
		})

		result[groupId] = group
	}

	return lo.Values(result), nil
}

func (manager *Manager) Resolve(ctx context.Context, values ...deps.Dependency) error {
	validStatuses := []deps.Status{
		deps.OK,
		deps.Failed,
	}

	values = lo.UniqBy(
		lo.Filter(values, func(item deps.Dependency, _ int) bool {
			return lo.Contains(validStatuses, item.Status)
		}),
		func(item deps.Dependency) string {
			return item.ID
		},
	)

	status2deps := lo.GroupBy(values, func(dep deps.Dependency) deps.Status {
		return dep.Status
	})

	sql := `
		with cte as (
			update dependencies
				set status = $1
				where id = any ($2) and status = $3
				returning *)
		update groups
		set pending = pending - 1
		where id in
			  (select distinct gd.group_id
			   from cte
						join group_dependencies gd
							 on gd.dependency_id = cte.id)`

	tx, err := manager.client.Begin(ctx)
	if err != nil {
		return fmt.Errorf("failed to begin a transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	for status, _deps := range status2deps {
		ids := lo.Map(_deps, func(dep deps.Dependency, _ int) string {
			return dep.ID
		})
		if _, err := tx.Exec(ctx, sql, status, ids, deps.Pending); err != nil {
			return fmt.Errorf("failed to resolve dependencies (status=%s): %w", status, err)
		}
	}
	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("failed to commit the transaction: %w", err)
	}

	return nil
}

func (manager *Manager) NewDependencies(ctx context.Context, n int) ([]deps.Dependency, error) {
	ids := make([]string, 0, n)
	for i := 0; i < n; i++ {
		ids = append(ids, manager.generateNewID())
	}

	newDependencies := lo.Map(ids, func(id string, _ int) deps.Dependency {
		return deps.Dependency{
			ID:     id,
			Status: deps.Pending,
		}
	})

	_, err := manager.client.CopyFrom(ctx,
		pgx.Identifier{"dependencies"},
		[]string{"id", "status"},
		pgx.CopyFromRows(
			lo.Map(newDependencies, func(dep deps.Dependency, _ int) []any {
				return []any{dep.ID, dep.Status}
			})),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to insert dependencies: %w", err)
	}

	return newDependencies, nil
}

func (manager *Manager) NewGroup(ctx context.Context, ids ...string) (groupId string, err error) {
	groupId = manager.generateNewID()

	tx, err := manager.client.Begin(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to begin a transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	// Initializing the group
	groupCreationQuery := `
		INSERT INTO groups (id, pending, status) 
		VALUES ($1, 0, $2)
	`

	_, err = tx.Exec(ctx, groupCreationQuery, groupId, GroupInitializingStatus)
	if err != nil {
		return "", fmt.Errorf("failed to initialize group: %w", err)
	}

	// Initializing group dependencies
	_, err = tx.CopyFrom(ctx,
		pgx.Identifier{"group_dependencies"},
		[]string{"dependency_id", "group_id"},
		pgx.CopyFromRows(
			lo.Map(ids, func(id string, _ int) []any {
				return []any{id, groupId}
			})),
	)
	if err != nil {
		return "", fmt.Errorf("failed to insert group's dependencies: %w", err)
	}

	// Updating the group status
	groupStatusUpdateQuery := `
		UPDATE groups
		SET pending = (
		    SELECT COUNT(*) FROM group_dependencies gd
				JOIN dependencies d ON d.id = gd.dependency_id
				WHERE d.status = $2 and gd.group_id = $3
		), status = $1
		WHERE id = $3
		
	`

	_, err = tx.Exec(ctx,
		groupStatusUpdateQuery,
		GroupWaitingStatus,
		deps.Pending,
		groupId,
	)
	if err != nil {
		return "", fmt.Errorf("failed to update a group's status: %w", err)
	}

	if err := tx.Commit(ctx); err != nil {
		return "", fmt.Errorf("failed to commit: %w", err)
	}

	return groupId, nil
}

func (manager *Manager) ReadyGroups(ctx context.Context) (<-chan string, error) {
	channel := make(chan string, 256)

	go manager.pollReadyGroups(ctx, channel)

	go func(ctx context.Context) {
		<-ctx.Done()
		close(channel)
	}(ctx)

	return channel, nil
}

func (manager *Manager) pollReadyGroups(ctx context.Context, channel chan<- string) {
	ticker := time.NewTicker(manager.config.PollingInterval)

poller:
	for {
		select {
		case <-ctx.Done():
			break poller
		case <-ticker.C:
			groups, err := manager.loadReadyGroups(ctx, manager.config.PollingBatchSize)
			if err != nil {
				slog.Error("failed to load ready groups",
					slog.String("error", err.Error()))
				continue
			}

		groupsIterator:
			for _, group := range groups {
				select {
				case <-ctx.Done():
					break groupsIterator
				case channel <- group:
				}
			}
		}
	}
}

func (manager *Manager) loadReadyGroups(ctx context.Context, limit int) ([]string, error) {
	sql := `
		update groups
		set status = $1
		where id in (select id from groups 
		                       where status = $2 and pending = 0
		                       limit $3)
		returning id
	`

	rows, err := manager.client.Query(ctx, sql,
		GroupCollectedStatus,
		GroupWaitingStatus,
		limit)
	if err != nil {
		return nil, fmt.Errorf("failed to query: %w", err)
	}

	result := make([]string, 0, limit)

	for rows.Next() {
		var id string

		if err := rows.Scan(&id); err != nil {
			return nil, fmt.Errorf("failed to scane: %w", err)
		}
		result = append(result, id)
	}

	return result, nil
}

func (manager *Manager) generateNewID() string {
	return strings.ReplaceAll(uuid.New().String(), "-", "")
}

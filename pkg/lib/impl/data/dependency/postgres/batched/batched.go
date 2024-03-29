package batched

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"github.com/ischenkx/kantoku/pkg/common/dependency"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/samber/lo"
	"log/slog"
	"strings"
	"time"
)

var _ dependency.Manager = (*Manager)(nil)

type Manager struct {
	Client *pgxpool.Pool
	Config Config
	Logger *slog.Logger
}

func (manager *Manager) LoadDependencies(ctx context.Context, ids ...string) ([]dependency.Dependency, error) {
	sql := `select id, status from dependencies where id in ($1)`

	rows, err := manager.Client.Query(ctx, sql, ids)
	if err != nil {
		return nil, fmt.Errorf("failed to query: %w", err)
	}

	result := make([]dependency.Dependency, 0, len(ids))

	for rows.Next() {
		var dep dependency.Dependency
		if err := rows.Scan(&dep.ID, &dep.Status); err != nil {
			return nil, fmt.Errorf("failed to scan: %w", err)
		}
		result = append(result, dep)
	}

	return result, nil
}

func (manager *Manager) LoadGroups(ctx context.Context, ids ...string) ([]dependency.Group, error) {
	sql := `
			select group_id, dependency_id, d.status
			from group_dependencies gd
			join dependencies d on d.id = gd.dependency_id
			where group_id in ($1)
	`

	rows, err := manager.Client.Query(ctx, sql, ids)
	if err != nil {
		return nil, fmt.Errorf("failed to query: %w", err)
	}

	result := make(map[string]dependency.Group, len(ids))

	for rows.Next() {
		var groupId, dependencyId, status string

		if err := rows.Scan(&groupId, &dependencyId, &status); err != nil {
			return nil, fmt.Errorf("failed to scane: %w", err)
		}
		if _, ok := result[groupId]; !ok {
			result[groupId] = dependency.Group{ID: groupId}
		}
		group := result[groupId]

		group.Dependencies = append(group.Dependencies, dependency.Dependency{
			ID:     dependencyId,
			Status: dependency.Status(status),
		})

		result[groupId] = group
	}

	return lo.Values(result), nil
}

func (manager *Manager) Resolve(ctx context.Context, values ...dependency.Dependency) error {
	validStatuses := []dependency.Status{
		dependency.OK,
		dependency.Failed,
	}

	values = lo.UniqBy(
		lo.Filter(values, func(item dependency.Dependency, _ int) bool {
			return lo.Contains(validStatuses, item.Status)
		}),
		func(item dependency.Dependency) string {
			return item.ID
		},
	)

	status2deps := lo.GroupBy(values, func(dep dependency.Dependency) dependency.Status {
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

	tx, err := manager.Client.Begin(ctx)
	if err != nil {
		return fmt.Errorf("failed to begin a transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	for status, _deps := range status2deps {
		ids := lo.Map(_deps, func(dep dependency.Dependency, _ int) string {
			return dep.ID
		})
		if _, err := tx.Exec(ctx, sql, status, ids, dependency.Pending); err != nil {
			return fmt.Errorf("failed to resolve dependencies (status=%s): %w", status, err)
		}
	}
	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("failed to commit the transaction: %w", err)
	}

	return nil
}

func (manager *Manager) NewDependencies(ctx context.Context, n int) ([]dependency.Dependency, error) {
	ids := make([]string, 0, n)
	for i := 0; i < n; i++ {
		ids = append(ids, manager.generateNewID())
	}

	newDependencies := lo.Map(ids, func(id string, _ int) dependency.Dependency {
		return dependency.Dependency{
			ID:     id,
			Status: dependency.Pending,
		}
	})

	_, err := manager.Client.CopyFrom(ctx,
		pgx.Identifier{"dependencies"},
		[]string{"id", "status"},
		pgx.CopyFromRows(
			lo.Map(newDependencies, func(dep dependency.Dependency, _ int) []any {
				return []any{dep.ID, dep.Status}
			})),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to insert dependencies: %w", err)
	}

	return newDependencies, nil
}

func (manager *Manager) NewGroup(ctx context.Context) (groupId string, err error) {
	groupId = manager.generateNewID()

	tx, err := manager.Client.Begin(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to begin a transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	// Initializing the group
	groupCreationQuery := `
		INSERT INTO groups (id, pending, status) 
		VALUES ($1, 0, $2)
	`

	_, err = tx.Exec(ctx, groupCreationQuery, groupId, GroupCreatedStatus)
	if err != nil {
		return "", fmt.Errorf("failed to initialize group: %w", err)
	}

	if err := tx.Commit(ctx); err != nil {
		return "", fmt.Errorf("failed to commit: %w", err)
	}

	return groupId, nil
}

func (manager *Manager) InitializeGroup(ctx context.Context, groupId string, ids ...string) error {
	tx, err := manager.Client.Begin(ctx)
	if err != nil {
		return fmt.Errorf("failed to begin a transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	result, err := tx.Exec(ctx, `
		UPDATE groups 
		SET status = $1
		WHERE id = $2 AND status = $3
	`, GroupInitializingStatus, groupId, GroupCreatedStatus)
	if err != nil {
		return fmt.Errorf("failed to update the group status: %w", err)
	}
	if result.RowsAffected() == 0 {
		return fmt.Errorf("group is already initialized or doesn't exist")
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
		return fmt.Errorf("failed to insert group's dependencies: %w", err)
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
		dependency.Pending,
		groupId,
	)
	if err != nil {
		return fmt.Errorf("failed to update a group's status: %w", err)
	}

	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("failed to commit: %w", err)
	}

	return nil
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
	ticker := time.NewTicker(manager.Config.PollingInterval)

poller:
	for {
		select {
		case <-ctx.Done():
			break poller
		case <-ticker.C:
			groups, err := manager.loadReadyGroups(ctx, manager.Config.PollingBatchSize)
			if err != nil {
				manager.Logger.Error("failed to load ready groups",
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

	rows, err := manager.Client.Query(ctx, sql,
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

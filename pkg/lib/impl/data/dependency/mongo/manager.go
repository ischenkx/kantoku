package mongorep

import (
	"context"
	"fmt"
	"github.com/ischenkx/kantoku/pkg/common/dependency"
	"github.com/ischenkx/kantoku/pkg/common/transport/broker"
	"github.com/ischenkx/kantoku/pkg/core/event"
	"github.com/samber/lo"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log/slog"
	"time"
)

// Groups and Dependencies are stored in one collection and are differentiated by the 'kind' field.
type Manager struct {
	Collection    *mongo.Collection
	Broker        *event.Broker
	Logger        *slog.Logger
	EventName     string
	ConsumerGroup string
}

func (manager *Manager) LoadDependencies(ctx context.Context, ids ...string) ([]dependency.Dependency, error) {
	filter := bson.M{
		"_id": bson.M{
			"$in": ids,
		},
		"kind": DependencyKind,
	}

	cursor, err := manager.Collection.Find(ctx, filter)
	if err != nil {
		return nil, fmt.Errorf("failed to find: %w", err)
	}
	defer cursor.Close(ctx)

	var models []Dependency
	if err := cursor.All(ctx, &models); err != nil {
		return nil, fmt.Errorf("failed to decode: %w", err)
	}

	deps := lo.Map(models, func(dep Dependency, _ int) dependency.Dependency {
		return dependency.Dependency{
			ID:     dep.ID,
			Status: dependency.Status(dep.Status),
		}
	})

	return deps, nil
}

func (manager *Manager) LoadGroups(ctx context.Context, ids ...string) ([]dependency.Group, error) {
	shallowGroups, err := manager.LoadShallowGroups(ctx, ids...)
	if err != nil {
		return nil, fmt.Errorf("failed to load shallow groups: %w", err)
	}

	dependencyIds := lo.Uniq(
		lo.FlatMap(shallowGroups, func(group Group, _ int) []string {
			groupResources := make([]string, 0, len(group.Pending)+len(group.Ready))
			groupResources = append(groupResources, lo.Keys(group.Pending)...)
			groupResources = append(groupResources, lo.Keys(group.Ready)...)

			return groupResources
		}),
	)

	deps, err := manager.LoadDependencies(ctx, dependencyIds...)
	if err != nil {
		return nil, fmt.Errorf("failed to load dependencies: %w", err)
	}

	index := lo.SliceToMap(deps, func(dep dependency.Dependency) (string, dependency.Dependency) {
		return dep.ID, dep
	})

	groups := lo.Map(shallowGroups, func(shallowGroup Group, _ int) dependency.Group {
		dependencies := make([]dependency.Dependency, 0, len(shallowGroup.Pending)+len(shallowGroup.Ready))
		for id := range shallowGroup.Pending {
			dependencies = append(dependencies, index[id])
		}
		for id := range shallowGroup.Ready {
			dependencies = append(dependencies, index[id])
		}

		return dependency.Group{
			ID:           shallowGroup.ID,
			Dependencies: dependencies,
		}
	})

	return groups, nil
}

func (manager *Manager) Resolve(ctx context.Context, values ...dependency.Dependency) error {
	values = lo.UniqBy(values, func(dep dependency.Dependency) string {
		return dep.ID
	})
	values = lo.Filter(values, func(dep dependency.Dependency, _ int) bool {
		return dep.Status != dependency.Pending
	})
	if len(values) == 0 {
		return nil
	}

	ids := lo.Map(values, func(dep dependency.Dependency, _ int) string {
		return dep.ID
	})

	now := time.Now().UnixNano()

	// Updating dependencies
	_, err := manager.Collection.UpdateMany(ctx,
		bson.M{
			"kind":   DependencyKind,
			"status": dependency.Pending,
			"_id": bson.M{
				"$in": ids,
			},
		},
		bson.M{
			"$set": bson.M{
				"status": bson.M{
					"$switch": lo.Map(values, func(dep dependency.Dependency, _ int) bson.M {
						return bson.M{
							"case": bson.M{"_id": bson.M{"$eq": dep.ID}},
							"then": dep.Status,
						}
					}),
				},
				"groups_processed": false,
				"updated_at":       now,
			},
		},
	)
	if err != nil {
		return fmt.Errorf("failed to update dependencies: %w", err)
	}

	// Updating groups

	// TODO:
	// 1. Update the dependencies and set substatus to groups-resolving
	// 2. Update groups
	// 3. Pull ready groups and emit them
	// 4. Update dependencies

	//TODO implement me
	panic("implement me")
}

func (manager *Manager) NewDependencies(ctx context.Context, n int) ([]dependency.Dependency, error) {
	if n <= 0 {
		return nil, nil
	}

	toBeInserted := lo.RepeatBy(n, func(_ int) any {
		return Dependency{
			Doc: Doc{
				ContextID: "",
				Kind:      DependencyKind,
				UpdatedAt: time.Now().UnixNano(),
			},
			Status: dependency.Pending,
		}
	})

	result, err := manager.Collection.InsertMany(ctx, toBeInserted, options.InsertMany())
	if err != nil {
		return nil, fmt.Errorf("failed to insert: %w", err)
	}

	deps := make([]dependency.Dependency, 0, len(result.InsertedIDs))
	for _, rawId := range result.InsertedIDs {
		id, isString := rawId.(string)
		if !isString {
			return nil, fmt.Errorf("failed to cast raw id to string")
		}

		deps = append(deps, dependency.Dependency{
			ID:     id,
			Status: dependency.Pending,
		})
	}

	return deps, nil
}

func (manager *Manager) NewGroup(ctx context.Context) (groupId string, err error) {
	result, err := manager.Collection.InsertOne(ctx, Group{
		Doc: Doc{
			ContextID: "",
			Kind:      GroupKind,
			UpdatedAt: time.Now().UnixNano(),
		},
		Status: GroupCreated,
	})
	if err != nil {
		return "", fmt.Errorf("failed to insert: %w", err)
	}

	id, ok := result.InsertedID.(string)
	if !ok {
		return "", fmt.Errorf("failed to cast raw id to string")
	}

	return id, nil
}

func (manager *Manager) InitializeGroup(ctx context.Context, groupId string, dependencyIds ...string) error {
	pending := Set[string]{}
	pending.Add(dependencyIds...)

	result, err := manager.Collection.UpdateOne(
		ctx,
		bson.M{
			"kind":   GroupKind,
			"status": GroupCreated,
			"_id":    groupId,
		},
		bson.M{
			"$set": bson.M{
				"pending": pending,
				"ready":   Set[string]{},
				"status":  GroupInitializing,
			},
		},
	)
	if err != nil {
		return fmt.Errorf("failed to update (initializing): %w", err)
	}
	if result.MatchedCount == 0 {
		return fmt.Errorf("new group not found")
	}

	dependencies, err := manager.LoadDependencies(ctx, dependencyIds...)
	if err != nil {
		return fmt.Errorf("failed to load dependencies: %w", err)
	}

	unset := bson.M{}
	set := bson.M{
		"status": GroupInitialized,
	}
	for _, dep := range dependencies {
		if dep.Status == dependency.Pending {
			continue
		}

		unset[fmt.Sprintf("pending.%s", dep.ID)] = ""
		set[fmt.Sprintf("ready.%s", dep.ID)] = nil
	}

	result, err = manager.Collection.UpdateOne(
		ctx,
		bson.M{
			"kind":   GroupKind,
			"status": GroupInitializing,
			"_id":    groupId,
		},
		bson.M{
			"$unset": unset,
			"$set":   set,
		},
	)
	if err != nil {
		return fmt.Errorf("failed to update (initialized): %w", err)
	}
	if result.MatchedCount == 0 {
		return fmt.Errorf("initializing group not found")
	}

	return nil
}

func (manager *Manager) ReadyGroups(ctx context.Context) (<-chan string, error) {
	channel, err := manager.Broker.Consume(ctx, broker.TopicsInfo{
		Group: manager.ConsumerGroup,
		Topics: []string{
			manager.EventName,
		},
	})
	if err != nil {
		return nil, err
	}

	outputs := make(chan string, 1024)

	go func(ctx context.Context, outputs chan<- string) {
		defer close(outputs)

		queue.Processor[event.Event]{
			Handler: func(ctx context.Context, ev event.Event) error {
				select {
				case <-ctx.Done():
					return fmt.Errorf("context done")
				case outputs <- string(ev.Data):
					return nil
				}
			},
			ErrorHandler: func(ctx context.Context, ev event.Event, err error) {
				manager.Logger.Error("failed to convert the event to id:", err)
			},
		}.Process(ctx, channel)
	}(ctx, outputs)

	return outputs, nil
}

func (manager *Manager) LoadShallowGroups(ctx context.Context, ids ...string) ([]Group, error) {
	filter := bson.M{
		"_id": bson.M{
			"$in": ids,
		},
		"kind": GroupKind,
	}

	cursor, err := manager.Collection.Find(ctx, filter)
	if err != nil {
		return nil, fmt.Errorf("failed to find: %w", err)
	}
	defer cursor.Close(ctx)

	var models []Group
	if err := cursor.All(ctx, &models); err != nil {
		return nil, fmt.Errorf("failed to decode: %w", err)
	}

	return models, nil
}

func (manager *Manager) syncGroups(ctx context.Context, dependencies ...string) {

}

func (manager *Manager) sendReadyGroup(ctx context.Context, id string) error {
	return manager.Broker.Send(ctx, event.New(manager.EventName, []byte(id)))
}

func (manager *Manager) sendReadyGroupQuietly(ctx context.Context, id string) {
	if err := manager.sendReadyGroup(ctx, id); err != nil {
		manager.Logger.Error("failed to send a ready group:", err)
	}
}

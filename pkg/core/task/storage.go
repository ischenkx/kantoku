package task

import (
	"context"
	"fmt"
	"github.com/ischenkx/kantoku/pkg/common/data/codec"
	"github.com/ischenkx/kantoku/pkg/common/data/storage"
	"go.mongodb.org/mongo-driver/bson"
	"reflect"
)

type Storage interface {
	storage.Storage

	Insert(ctx context.Context, tasks []Task) error
	Delete(ctx context.Context, ids []string) error
	ByIDs(ctx context.Context, ids []string) ([]Task, error)
	UpdateByIDs(ctx context.Context, ids []string, properties map[string]any) error
	GetWithProperties(ctx context.Context, propertiesToValues map[string][]any) ([]Task, error)
	UpdateWithProperties(ctx context.Context, propertiesToValues map[string][]any, newProperties map[string]any) (updatedDocs int, err error)
}

var _ Storage = (*MongoStorage)(nil)

type MongoStorage struct {
	BaseStorage *storage.MongoStorage
	Codec       codec.Codec[Task, map[string]any]
}

func (ms *MongoStorage) Settings(ctx context.Context) (storage.Settings, error) {
	return ms.BaseStorage.Settings(ctx)
}

func (ms *MongoStorage) Exec(ctx context.Context, command storage.Command) ([]storage.Document, error) {
	return ms.BaseStorage.Exec(ctx, command)
}

func (ms *MongoStorage) Insert(ctx context.Context, tasks []Task) error {
	encoded := make([]storage.Document, 0, len(tasks))
	for _, task := range tasks {
		encodedTask, err := ms.Codec.Encode(task)
		if err != nil {
			return fmt.Errorf("failed to encode task: %w", err)
		}
		encoded = append(encoded, encodedTask)
	}
	res, err := ms.Exec(ctx, storage.Command{
		Operation: "insert",
		Params: []storage.Param{
			{
				Name:  "documents",
				Value: encoded,
			},
		},
	})

	_ = res

	if err != nil {
		return err
	}

	return nil
}

func (ms *MongoStorage) Delete(ctx context.Context, ids []string) error {
	_, err := ms.Exec(ctx, storage.Command{
		Operation: "delete",
		Params: []storage.Param{
			{
				Name: "deletes",
				Value: []map[string]any{
					{
						"q": map[string]any{
							"id": map[string]any{
								"$in": ids,
							},
						},
						"limit": 0,
					},
				},
			},
		},
	})

	if err != nil {
		return err
	}

	return nil
}

func (ms *MongoStorage) ByIDs(ctx context.Context, ids []string) ([]Task, error) {
	docs, err := ms.Exec(ctx, storage.Command{
		Operation: "find",
		Params: []storage.Param{
			{
				Name: "filter",
				Value: map[string]any{
					"id": map[string]any{
						"$in": ids,
					},
				},
			},
		},
	})
	if err != nil {
		return nil, err
	}

	if len(docs) == 0 {
		return nil, fmt.Errorf("no docs returned")
	}

	cursorDoc := docs[0]
	cursor, ok := cursorDoc["cursor"]
	if !ok {
		return nil, fmt.Errorf("no property 'cursor' in the first doc")
	}

	cursorObject, ok := cursor.(map[string]any)
	if !ok {
		return nil, fmt.Errorf("cursor has the wrong type")
	}

	firstBatch, ok := cursorObject["firstBatch"]
	if !ok {
		return nil, fmt.Errorf("no property 'firstBatch' in the cursor")
	}

	firstBatchArray, ok := firstBatch.(bson.A)
	if !ok {
		return nil, fmt.Errorf("firstBatch has the wrong type: %s", reflect.TypeOf(firstBatchArray))
	}

	var result []Task
	for _, doc := range firstBatchArray {
		docObject, ok := doc.(map[string]interface{})
		if !ok {
			return nil, fmt.Errorf("doc is not an object")
		}

		t, err := ms.Codec.Decode(docObject)
		if err != nil {
			return nil, fmt.Errorf("failed to decode task: %w", err)
		}

		result = append(result, t)
	}

	return result, nil
}

func (ms *MongoStorage) UpdateByIDs(ctx context.Context, ids []string, properties map[string]any) error {
	_, err := ms.Exec(ctx, storage.Command{
		Operation: "update",
		Params: []storage.Param{
			{
				Name: "updates",
				Value: []any{
					map[string]any{
						"q": map[string]any{
							"id": map[string]any{
								"$in": ids,
							},
						},
						"u": map[string]any{
							"$set": properties,
						},
					},
				},
			},
		},
	})
	if err != nil {
		return err
	}

	return nil
}

func (ms *MongoStorage) GetWithProperties(ctx context.Context, propertiesToValues map[string][]any) ([]Task, error) {
	filter := map[string]any{}
	for key, value := range propertiesToValues {
		filter[key] = map[string]any{
			"$in": value,
		}
	}

	docs, err := ms.Exec(ctx, storage.Command{
		Operation: "find",
		Params: []storage.Param{
			{
				Name:  "filter",
				Value: filter,
			},
		},
	})
	if err != nil {
		return nil, err
	}

	if len(docs) == 0 {
		return nil, fmt.Errorf("no docs returned")
	}

	cursorDoc := docs[0]
	cursor, ok := cursorDoc["cursor"]
	if !ok {
		return nil, fmt.Errorf("no property 'cursor' in the first doc")
	}

	cursorObject, ok := cursor.(map[string]any)
	if !ok {
		return nil, fmt.Errorf("cursor has the wrong type")
	}

	firstBatch, ok := cursorObject["firstBatch"]
	if !ok {
		return nil, fmt.Errorf("no property 'firstBatch' in the cursor")
	}

	firstBatchArray, ok := firstBatch.([]any)
	if !ok {
		return nil, fmt.Errorf("firstBatch has the wrong type")
	}

	var result []Task
	for _, doc := range firstBatchArray {
		docObject, ok := doc.(map[string]interface{})
		if !ok {
			return nil, fmt.Errorf("doc is not an object")
		}

		t, err := ms.Codec.Decode(docObject)
		if err != nil {
			return nil, fmt.Errorf("failed to decode task: %w", err)
		}

		result = append(result, t)
	}

	return result, nil
}

func (ms *MongoStorage) UpdateWithProperties(ctx context.Context, propertiesToValues map[string][]any, newProperties map[string]any) (int, error) {
	filter := map[string]any{}
	for key, value := range propertiesToValues {
		filter[key] = map[string]any{
			"$in": value,
		}
	}

	docs, err := ms.Exec(ctx, storage.Command{
		Operation: "update",
		Params: []storage.Param{
			{
				Name: "updates",
				Value: []any{
					map[string]any{
						"q": filter,
						"u": map[string]any{
							"$set": newProperties,
						},
					},
				},
			},
		},
	})
	if err != nil {
		return 0, err
	}

	if len(docs) == 0 {
		return 0, fmt.Errorf("expected 1 document to be returned, got 0")
	}

	doc := docs[0]
	rawOk, exists := doc["ok"]
	if !exists {
		return 0, fmt.Errorf("expected 'ok' in the doc")
	}

	if ok, asserted := rawOk.(bool); asserted && !ok {
		return 0, fmt.Errorf("operation failed")
	}

	rawModified, ok := doc["n"]
	if !ok {
		return 0, fmt.Errorf("expected 'n' in the doc")
	}

	modified, ok := getInt(rawModified)
	if !ok {
		return 0, fmt.Errorf("failed to get amount of modified docs")
	}

	return modified, nil
}

func getInt(val any) (int, bool) {
	res1, isFloat64 := val.(float64)
	if isFloat64 {
		return int(res1), true
	}

	res2, isFloat32 := val.(float32)
	if isFloat32 {
		return int(res2), true
	}

	res3, isInt64 := val.(int64)
	if isInt64 {
		return int(res3), true
	}

	res4, isInt32 := val.(int32)
	if isInt32 {
		return int(res4), true
	}

	return 0, false
}

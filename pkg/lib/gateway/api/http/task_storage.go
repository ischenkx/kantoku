package http

import (
	"context"
	"fmt"
	"github.com/ischenkx/kantoku/pkg/common/data/storage"
	"github.com/ischenkx/kantoku/pkg/core/task"
	"github.com/ischenkx/kantoku/pkg/lib/gateway/api/http/oas"
	"github.com/samber/lo"
	"net/http"
)

var _ task.Storage = (*taskStorage)(nil)

type taskStorage struct {
	httpClient oas.ClientWithResponsesInterface
}

func (t taskStorage) Settings(ctx context.Context) (storage.Settings, error) {
	resp, err := t.httpClient.PostTasksStorageSettingsWithResponse(ctx)
	if err != nil {
		return storage.Settings{}, err
	}

	switch resp.StatusCode() {
	case http.StatusOK:
		return storage.Settings{
			Type: resp.JSON200.Type,
			Meta: resp.JSON200.Meta,
		}, nil
	case http.StatusInternalServerError:
		return storage.Settings{}, fmt.Errorf("server failure: %s", resp.JSON500.Message)
	default:
		return storage.Settings{}, fmt.Errorf("unexpected response code: %d", resp.StatusCode())
	}
}

func (t taskStorage) Exec(ctx context.Context, command storage.Command) ([]storage.Document, error) {
	res, err := t.httpClient.PostTasksStorageExecWithResponse(ctx, oas.TaskStorageCommand{
		Meta:      command.Meta,
		Operation: command.Operation,
		Params: lo.Map(command.Params, func(param storage.Param, _ int) oas.TaskStorageCommandParam {
			return oas.TaskStorageCommandParam{
				Name:  param.Name,
				Value: param.Value,
			}
		}),
	})
	if err != nil {
		return nil, err
	}

	switch res.StatusCode() {
	case http.StatusOK:
		return lo.Map(*res.JSON200, func(doc map[string]any, _ int) storage.Document {
			return doc
		}), nil
	case http.StatusInternalServerError:
		return nil, fmt.Errorf("server failure: %s", res.JSON500.Message)
	default:
		return nil, fmt.Errorf("unexpected response code: %d", res.StatusCode())
	}
}

func (t taskStorage) Insert(ctx context.Context, tasks []task.Task) error {
	res, err := t.httpClient.PostTasksStorageInsertWithResponse(ctx, lo.Map(tasks, func(t task.Task, _ int) oas.Task {
		return TaskToDto(t)
	}))
	if err != nil {
		return err
	}

	switch res.StatusCode() {
	case http.StatusOK:
		return nil
	case http.StatusInternalServerError:
		return fmt.Errorf("server failure: %s", res.JSON500.Message)
	default:
		return fmt.Errorf("unexpected response code: %d", res.StatusCode())
	}
}

func (t taskStorage) Delete(ctx context.Context, ids []string) error {
	res, err := t.httpClient.PostTasksStorageDeleteWithResponse(ctx, ids)
	if err != nil {
		return err
	}

	switch res.StatusCode() {
	case http.StatusOK:
		return nil
	case http.StatusInternalServerError:
		return fmt.Errorf("server failure: %s", res.JSON500.Message)
	default:
		return fmt.Errorf("unexpected response code: %d", res.StatusCode())
	}
}

func (t taskStorage) ByIDs(ctx context.Context, ids []string) ([]task.Task, error) {
	res, err := t.httpClient.PostTasksStorageGetByIdsWithResponse(ctx, ids)
	if err != nil {
		return nil, err
	}

	switch res.StatusCode() {
	case http.StatusOK:
		return lo.Map(*res.JSON200, func(t oas.Task, _ int) task.Task {
			return task.Task{
				Inputs:  t.Inputs,
				Outputs: t.Outputs,
				ID:      t.Id,
				Info:    t.Info,
			}
		}), nil
	case http.StatusInternalServerError:
		return nil, fmt.Errorf("server failure: %s", res.JSON500.Message)
	default:
		return nil, fmt.Errorf("unexpected response code: %d", res.StatusCode())
	}
}

func (t taskStorage) UpdateByIDs(ctx context.Context, ids []string, properties map[string]any) error {
	res, err := t.httpClient.PostTasksStorageUpdateByIdsWithResponse(ctx, oas.PostTasksStorageUpdateByIdsJSONRequestBody{
		Ids:        ids,
		Properties: properties,
	})
	if err != nil {
		return err
	}

	switch res.StatusCode() {
	case http.StatusOK:
		return nil
	case http.StatusInternalServerError:
		return fmt.Errorf("server failure: %s", res.JSON500.Message)
	default:
		return fmt.Errorf("unexpected response code: %d", res.StatusCode())
	}
}

func (t taskStorage) GetWithProperties(ctx context.Context, propertiesToValues map[string][]any) ([]task.Task, error) {
	res, err := t.httpClient.PostTasksStorageGetWithPropertiesWithResponse(ctx, oas.PostTasksStorageGetWithPropertiesJSONRequestBody{
		PropertiesToValues: propertiesToValues,
	})
	if err != nil {
		return nil, err
	}

	switch res.StatusCode() {
	case http.StatusOK:
		return lo.Map(*res.JSON200, func(t oas.Task, _ int) task.Task {
			return task.Task{
				Inputs:  t.Inputs,
				Outputs: t.Outputs,
				ID:      t.Id,
				Info:    t.Info,
			}
		}), nil
	case http.StatusInternalServerError:
		return nil, fmt.Errorf("server failure: %s", res.JSON500.Message)
	default:
		return nil, fmt.Errorf("unexpected response code: %d", res.StatusCode())
	}
}

func (t taskStorage) UpdateWithProperties(ctx context.Context, propertiesToValues map[string][]any, newProperties map[string]any) (int, error) {
	res, err := t.httpClient.PostTasksStorageUpdateWithPropertiesWithResponse(ctx, oas.PostTasksStorageUpdateWithPropertiesJSONRequestBody{
		NewProperties:      newProperties,
		PropertiesToValues: propertiesToValues,
	})
	if err != nil {
		return 0, err
	}

	switch res.StatusCode() {
	case http.StatusOK:
		return res.JSON200.Modified, nil
	case http.StatusInternalServerError:
		return 0, fmt.Errorf("server failure: %s", res.JSON500.Message)
	default:
		return 0, fmt.Errorf("unexpected response code: %d", res.StatusCode())
	}
}

package client

import (
	"context"
	"errors"
	"fmt"
	"github.com/ischenkx/kantoku/pkg/extensions/api/http/converters"
	"github.com/ischenkx/kantoku/pkg/extensions/api/http/oas"
	"github.com/ischenkx/kantoku/pkg/system/kernel/task"
	"github.com/samber/lo"
	"net/http"
)

var _ task.Storage = (*taskStorage)(nil)

type taskStorage struct {
	httpClient oas.ClientWithResponsesInterface
}

func (storage taskStorage) Create(_ context.Context, _ task.Task) (task.Task, error) {
	return task.Task{}, errors.New("not supported by an http client")
}

func (storage taskStorage) Delete(ctx context.Context, ids ...string) error {
	return errors.New("not supported by an http client")
}

func (storage taskStorage) Load(ctx context.Context, ids ...string) ([]task.Task, error) {
	res, err := storage.httpClient.PostTasksLoadWithResponse(ctx, ids)
	if err != nil {
		return nil, fmt.Errorf("failed to make an http request: %w", err)
	}

	code := res.StatusCode()

	switch code {
	case http.StatusOK:
		return lo.Map(*res.JSON200, func(t oas.Task, _ int) task.Task {
			return task.Task{
				Inputs:     t.Inputs,
				Outputs:    t.Outputs,
				Properties: converters.DtoToProperties(t.Properties),
				ID:         t.Id,
			}
		}), nil
	case http.StatusInternalServerError:
		return nil, fmt.Errorf("server failure: %s", res.JSON500.Message)
	default:
		return nil, fmt.Errorf("unexpected response code: %d", code)
	}
}

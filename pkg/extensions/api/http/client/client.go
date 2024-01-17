package client

import (
	"context"
	"fmt"
	"github.com/ischenkx/kantoku/pkg/common/data/record"
	"github.com/ischenkx/kantoku/pkg/extensions/api/http/converters"
	"github.com/ischenkx/kantoku/pkg/extensions/api/http/oas"
	"github.com/ischenkx/kantoku/pkg/system"
	event "github.com/ischenkx/kantoku/pkg/system/kernel/event"
	"github.com/ischenkx/kantoku/pkg/system/kernel/resource"
	"github.com/ischenkx/kantoku/pkg/system/kernel/task"
	"net/http"
)

var _ system.AbstractSystem = (*Client)(nil)

type Client struct {
	httpClient oas.ClientWithResponsesInterface
}

func New(clientInterface oas.ClientWithResponsesInterface) *Client {
	return &Client{
		httpClient: clientInterface,
	}
}

func (client *Client) Tasks() task.Storage {
	return taskStorage{httpClient: client.httpClient}
}

func (client *Client) Resources() resource.Storage {
	return resourceStorage{httpClient: client.httpClient}
}

func (client *Client) Events() event.Bus {
	return eventBus{}
}

func (client *Client) Info() record.Storage {
	return recordStorage{httpClient: client.httpClient}
}

func (client *Client) Spawn(ctx context.Context, initializers ...system.TaskInitializer) (*system.Task, error) {
	t := task.Task{}

	for _, initializer := range initializers {
		if initializer == nil {
			continue
		}

		initializer(&t)
	}

	res, err := client.httpClient.PostTasksSpawnWithResponse(ctx, oas.PostTasksSpawnJSONRequestBody{
		Inputs:     t.Inputs,
		Outputs:    t.Outputs,
		Properties: converters.PropertiesToDto(t.Properties),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to make an http request: %w", err)
	}

	code := res.StatusCode()

	switch code {
	case http.StatusOK:
		return client.Task(res.JSON200.Id), nil
	case http.StatusInternalServerError:
		return nil, fmt.Errorf("server failure: %s", res.JSON500.Message)
	default:
		return nil, fmt.Errorf("unexpected response code: %d", code)
	}
}

func (client *Client) Task(id string) *system.Task {
	return &system.Task{
		ID:     id,
		System: client,
	}
}

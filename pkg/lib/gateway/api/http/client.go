package http

import (
	"context"
	"fmt"
	"github.com/ischenkx/kantoku/pkg/core/event"
	"github.com/ischenkx/kantoku/pkg/core/resource"
	"github.com/ischenkx/kantoku/pkg/core/system"
	"github.com/ischenkx/kantoku/pkg/core/task"
	"github.com/ischenkx/kantoku/pkg/lib/gateway/api/http/oas"
	"net/http"
)

var _ system.AbstractSystem = (*Client)(nil)

type Client struct {
	httpClient oas.ClientWithResponsesInterface
}

func NewClient(clientInterface oas.ClientWithResponsesInterface) *Client {
	return &Client{
		httpClient: clientInterface,
	}
}

func (client *Client) Specifications() *SpecificationStorage {
	return &SpecificationStorage{client.httpClient}
}

func (client *Client) Types() *SpecificationStorage {
	return &SpecificationStorage{client.httpClient}
}

func (client *Client) Tasks() task.Storage {
	return taskStorage{client.httpClient}
}

func (client *Client) Resources() resource.Storage {
	return resourceStorage{httpClient: client.httpClient}
}

func (client *Client) Events() *event.Broker {
	return nil
}

func (client *Client) Spawn(ctx context.Context, t task.Task) (task.Task, error) {
	res, err := client.httpClient.PostTasksSpawnWithResponse(ctx, oas.PostTasksSpawnJSONRequestBody{
		Info:    t.Info,
		Inputs:  t.Inputs,
		Outputs: t.Outputs,
	})
	if err != nil {
		return t, fmt.Errorf("failed to make an http request: %w", err)
	}

	code := res.StatusCode()

	switch code {
	case http.StatusOK:
		return client.Task(ctx, res.JSON200.Id)
	case http.StatusInternalServerError:
		return t, fmt.Errorf("server failure: %s", res.JSON500.Message)
	default:
		return t, fmt.Errorf("unexpected response code: %d", code)
	}
}

func (client *Client) Task(ctx context.Context, id string) (task.Task, error) {
	ts, err := client.Tasks().ByIDs(ctx, []string{id})
	if err != nil {
		return task.Task{}, nil
	}

	if len(ts) == 0 {
		return task.Task{}, fmt.Errorf("not found")
	}

	return ts[0], nil
}

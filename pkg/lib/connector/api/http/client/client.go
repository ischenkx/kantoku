package client

import (
	"context"
	"fmt"
	"github.com/ischenkx/kantoku/pkg/common/data/record"
	recutil "github.com/ischenkx/kantoku/pkg/common/data/record/util"
	"github.com/ischenkx/kantoku/pkg/core/event"
	"github.com/ischenkx/kantoku/pkg/core/resource"
	"github.com/ischenkx/kantoku/pkg/core/system"
	"github.com/ischenkx/kantoku/pkg/core/task"
	"github.com/ischenkx/kantoku/pkg/lib/connector/api/http/oas"
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

func (client *Client) Tasks() record.Set[task.Task] {
	return recordStorage{httpClient: client.httpClient}
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
	t, err := recutil.Single(
		ctx,
		client.
			Tasks().
			Filter(record.R{"id": id}).
			Cursor().
			Iter(),
	)
	if err != nil {
		return task.Task{}, err
	}

	return t, nil
}

package kantokuhttp

import (
	"context"
	"fmt"
	"github.com/ischenkx/kantoku/pkg/core"
	"github.com/ischenkx/kantoku/pkg/lib/gateway/api/kantokuhttp/oas"
	"net/http"
)

var _ core.AbstractSystem = (*Client)(nil)

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

func (client *Client) Tasks() core.TaskDB {
	return taskStorage{client.httpClient}
}

func (client *Client) Resources() core.ResourceDB {
	return resourceStorage{httpClient: client.httpClient}
}

func (client *Client) Events() core.Broker {
	return nil
}

func (client *Client) Spawn(ctx context.Context, t core.Task) (core.Task, error) {
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

func (client *Client) Task(ctx context.Context, id string) (core.Task, error) {
	ts, err := client.Tasks().ByIDs(ctx, []string{id})
	if err != nil {
		return core.Task{}, nil
	}

	if len(ts) == 0 {
		return core.Task{}, fmt.Errorf("not found")
	}

	return ts[0], nil
}

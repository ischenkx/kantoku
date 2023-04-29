package simple

import (
	"context"
	"kantoku/common/data/pool"
	"time"
)

type Client struct {
	inputs  pool.Writer[Event]
	outputs pool.Reader[string]
}

func NewClient(inputs pool.Writer[Event], outputs pool.Reader[string]) *Client {
	return &Client{
		inputs:  inputs,
		outputs: outputs,
	}
}

func (c *Client) Schedule(ctx context.Context, at time.Time, event string) error {
	return c.inputs.Write(ctx, Event{
		ID:   event,
		When: at,
	})
}

func (c *Client) Events(ctx context.Context) (<-chan string, error) {
	return c.outputs.Read(ctx)
}

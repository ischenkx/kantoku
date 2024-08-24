package scraper

import (
	"github.com/ischenkx/kantoku/pkg/lib/tasks/fn"
	"github.com/ischenkx/kantoku/pkg/lib/tasks/fn/future"
)

type (
	ParsePageInput struct {
		Page future.Future[[]byte]
	}

	ParsePageOutput struct {
		Result future.Future[map[string]any]
	}

	ParsePage struct {
		fn.Function[ParsePage, ParsePageInput, ParsePageOutput]
	}
)

var (
	_ fn.AbstractFunction[ParsePageInput, ParsePageOutput] = (*ParsePage)(nil)
)

func (f ParsePage) Call(ctx *fn.Context, input ParsePageInput) (output ParsePageOutput, err error) {

	output.Result = future.FromValue(map[string]any{
		"title": "Title",
		"meta": map[string]any{
			"size": 42,
		},
	})

	return output, nil
}

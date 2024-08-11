package http

import (
	"bytes"
	"fmt"
	"github.com/ischenkx/kantoku/pkg/lib/tasks/fn_d"
	"github.com/ischenkx/kantoku/pkg/lib/tasks/fn_d/future"
	"io"
	"net/http"
	"net/url"
	"time"
)

type (
	DoInput struct {
		Url    future.Future[string]
		Method future.Future[string]
		Body   future.Future[string]
	}

	DoOutput struct {
		Code      future.Future[int]
		Body      future.Future[string]
		TimeStamp future.Future[time.Time]
	}

	Do struct {
		fn_d.Function[Do, DoInput, DoOutput]
	}
)

var (
	_ fn_d.AbstractFunction[DoInput, DoOutput] = (*Do)(nil)
)

func (task Do) Call(ctx *fn_d.Context, input DoInput) (output DoOutput, err error) {
	_url, err := url.Parse(input.Url.Value())
	if err != nil {
		return DoOutput{}, fmt.Errorf("failed to parse url: %w", err)
	}

	response, err := http.DefaultClient.Do(&http.Request{
		Method: input.Method.Value(),
		URL:    _url,
		Body:   io.NopCloser(bytes.NewReader([]byte(input.Body.Value()))),
	})
	if err != nil {
		return DoOutput{}, fmt.Errorf("failed to do http request: %w", err)
	}
	defer response.Body.Close()

	responseBody, err := io.ReadAll(response.Body)
	if err != nil {
		return DoOutput{}, fmt.Errorf("failed to read response body: %w", err)
	}

	return DoOutput{
		Code:      future.FromValue(response.StatusCode),
		Body:      future.FromValue(string(responseBody)),
		TimeStamp: future.FromValue(time.Now()),
	}, nil
}

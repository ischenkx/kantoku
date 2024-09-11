package kantoku

import (
	"fmt"
	"github.com/ischenkx/kantoku/pkg/core"
	kantokuHttp "github.com/ischenkx/kantoku/pkg/lib/gateway/api/kantokuhttp"
	"github.com/ischenkx/kantoku/pkg/lib/gateway/api/kantokuhttp/oas"
)

type Kantoku = core.AbstractSystem

func Connect(url string, options ...oas.ClientOption) (*kantokuHttp.Client, error) {
	oasClient, err := oas.NewClientWithResponses(url, options...)
	if err != nil {
		return nil, fmt.Errorf("failed to create the oas client: %w", err)
	}
	return kantokuHttp.NewClient(oasClient), nil
}

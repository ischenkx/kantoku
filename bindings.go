package kantoku

import (
	"fmt"
	kantokuHttp "github.com/ischenkx/kantoku/pkg/lib/gateway/api/http"
	"github.com/ischenkx/kantoku/pkg/lib/gateway/api/http/oas"
)

func Connect(url string, options ...oas.ClientOption) (*kantokuHttp.Client, error) {
	oasClient, err := oas.NewClientWithResponses(url, options...)
	if err != nil {
		return nil, fmt.Errorf("failed to create the oas client: %w", err)
	}
	return kantokuHttp.NewClient(oasClient), nil
}

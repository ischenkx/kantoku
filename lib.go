package kantoku

import (
	"fmt"
	"github.com/ischenkx/kantoku/pkg/core/system"
	kantohttp "github.com/ischenkx/kantoku/pkg/lib/connector/api/http"
	"github.com/ischenkx/kantoku/pkg/lib/connector/api/http/oas"
)

type Kantoku system.AbstractSystem

func Connect(url string, options ...oas.ClientOption) (Kantoku, error) {
	client, err := oas.NewClientWithResponses(url, options...)
	if err != nil {
		return nil, fmt.Errorf("failed to create an oas client: %w", err)
	}

	return kantohttp.NewClient(client), nil
}

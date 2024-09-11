package kantokuhttp

import (
	"context"
	"fmt"
	"github.com/ischenkx/kantoku/pkg/core"
	"github.com/ischenkx/kantoku/pkg/lib/gateway/api/kantokuhttp/oas"
	"github.com/samber/lo"
	"net/http"
)

var _ core.ResourceDB = (*resourceStorage)(nil)

type resourceStorage struct {
	httpClient oas.ClientWithResponsesInterface
}

func (storage resourceStorage) Load(ctx context.Context, ids ...string) ([]core.Resource, error) {
	res, err := storage.httpClient.PostResourcesLoadWithResponse(ctx, ids)
	if err != nil {
		return nil, fmt.Errorf("failed to make an http request: %w", err)
	}

	code := res.StatusCode()

	switch code {
	case http.StatusOK:
		return lo.Map(*res.JSON200, func(r oas.Resource, _ int) core.Resource {
			return core.Resource{
				Data:   []byte(r.Value),
				ID:     r.Id,
				Status: r.Status,
			}
		}), nil
	case http.StatusInternalServerError:
		return nil, fmt.Errorf("server failure: %s", res.JSON500.Message)
	default:
		return nil, fmt.Errorf("unexpected response code: %d", code)
	}
}

func (storage resourceStorage) Alloc(ctx context.Context, amount int) ([]string, error) {
	res, err := storage.httpClient.PostResourcesAllocateWithResponse(ctx,
		&oas.PostResourcesAllocateParams{
			Amount: amount,
		})
	if err != nil {
		return nil, fmt.Errorf("failed to make an http request: %w", err)
	}

	code := res.StatusCode()

	switch code {
	case http.StatusOK:
		return *res.JSON200, nil
	case http.StatusInternalServerError:
		return nil, fmt.Errorf("server failure: %s", res.JSON500.Message)
	default:
		return nil, fmt.Errorf("unexpected response code: %d", code)
	}
}

func (storage resourceStorage) Init(ctx context.Context, resources []core.Resource) error {
	res, err := storage.httpClient.PostResourcesInitializeWithResponse(ctx,
		lo.Map(resources, func(res core.Resource, _ int) oas.ResourceInitializer {
			return oas.ResourceInitializer{
				Id:    res.ID,
				Value: string(res.Data),
			}
		}))
	if err != nil {
		return nil
	}

	code := res.StatusCode()

	switch code {
	case http.StatusOK:
		return nil
	case http.StatusInternalServerError:
		return fmt.Errorf("server failure: %s", res.JSON500.Message)
	default:
		return fmt.Errorf("unexpected response code: %d", code)
	}
}

func (storage resourceStorage) Dealloc(ctx context.Context, ids []string) error {
	res, err := storage.httpClient.PostResourcesDeallocateWithResponse(ctx, ids)
	if err != nil {
		return fmt.Errorf("failed to make an http request: %w", err)
	}

	code := res.StatusCode()

	switch code {
	case http.StatusOK:
		return nil
	case http.StatusInternalServerError:
		return fmt.Errorf("server failure: %s", res.JSON500.Message)
	default:
		return fmt.Errorf("unexpected response code: %d", code)
	}
}

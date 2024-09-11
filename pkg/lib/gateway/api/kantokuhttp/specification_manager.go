package kantokuhttp

import (
	"context"
	"fmt"
	"github.com/ischenkx/kantoku/pkg/lib/gateway/api/kantokuhttp/oas"
	"github.com/ischenkx/kantoku/pkg/lib/tasks/specification"
	"github.com/ischenkx/kantoku/pkg/lib/tasks/specification/typing"
	"net/http"
)

type SpecificationStorage struct {
	httpClient oas.ClientWithResponsesInterface
}

func (storage *SpecificationStorage) Get(ctx context.Context, id string) (specification.Specification, error) {
	res, err := storage.httpClient.PostTasksSpecificationsGetWithResponse(
		ctx,
		oas.PostTasksSpecificationsGetJSONRequestBody{Id: id},
	)
	if err != nil {
		return specification.Specification{}, fmt.Errorf("request failed: %w", err)
	}

	switch res.StatusCode() {
	case http.StatusOK:
		return specificationFromModel(*res.JSON200), nil
	case http.StatusInternalServerError:
		return specification.Specification{}, fmt.Errorf("server failure: %s", res.JSON500.Message)
	default:
		return specification.Specification{}, fmt.Errorf("unexpected response code: %d", res.StatusCode())
	}
}

func (storage *SpecificationStorage) GetAll(ctx context.Context) ([]specification.Specification, error) {
	res, err := storage.httpClient.PostTasksSpecificationsGetAllWithResponse(ctx)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}

	switch res.StatusCode() {
	case http.StatusOK:
		models := *res.JSON200
		result := make([]specification.Specification, 0, len(models))

		for _, model := range models {
			result = append(result, specificationFromModel(model))
		}

		return result, nil
	case http.StatusInternalServerError:
		return nil, fmt.Errorf("server failure: %s", res.JSON500.Message)
	default:
		return nil, fmt.Errorf("unexpected response code: %d", res.StatusCode())
	}
}

func (storage *SpecificationStorage) Add(ctx context.Context, spec specification.Specification) error {
	res, err := storage.httpClient.PostTasksSpecificationsCreateWithResponse(ctx, specificationToModel(spec))
	if err != nil {
		return fmt.Errorf("request failed: %w", err)
	}

	switch res.StatusCode() {
	case http.StatusOK:
		return nil
	case http.StatusInternalServerError:
		return fmt.Errorf("server failure: %s", res.JSON500.Message)
	default:
		return fmt.Errorf("unexpected response code: %d", res.StatusCode())
	}
}

func (storage *SpecificationStorage) Remove(ctx context.Context, id string) error {
	res, err := storage.httpClient.PostTasksSpecificationsRemoveWithResponse(ctx,
		oas.PostTasksSpecificationsRemoveJSONRequestBody{Id: id},
	)
	if err != nil {
		return fmt.Errorf("request failed: %w", err)
	}

	switch res.StatusCode() {
	case http.StatusOK:
		return nil
	case http.StatusInternalServerError:
		return fmt.Errorf("server failure: %s", res.JSON500.Message)
	default:
		return fmt.Errorf("unexpected response code: %d", res.StatusCode())
	}
}

type TypeStorage struct {
	httpClient oas.ClientWithResponsesInterface
}

func (storage *TypeStorage) Get(ctx context.Context, id string) (specification.TypeWithID, error) {
	res, err := storage.httpClient.PostTasksSpecificationsTypesGetWithResponse(
		ctx,
		oas.PostTasksSpecificationsTypesGetJSONRequestBody{Id: id},
	)
	if err != nil {
		return specification.TypeWithID{}, fmt.Errorf("request failed: %w", err)
	}

	switch res.StatusCode() {
	case http.StatusOK:
		return specification.TypeWithID{
			ID:   res.JSON200.Id,
			Type: typeFromModel(res.JSON200.Type),
		}, nil
	case http.StatusInternalServerError:
		return specification.TypeWithID{}, fmt.Errorf("server failure: %s", res.JSON500.Message)
	default:
		return specification.TypeWithID{}, fmt.Errorf("unexpected response code: %d", res.StatusCode())
	}
}

func (storage *TypeStorage) GetAll(ctx context.Context) ([]specification.TypeWithID, error) {
	res, err := storage.httpClient.PostTasksSpecificationsTypesGetAllWithResponse(ctx)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}

	switch res.StatusCode() {
	case http.StatusOK:
		models := *res.JSON200
		result := make([]specification.TypeWithID, 0, len(models))

		for _, model := range models {
			result = append(result, specification.TypeWithID{
				ID:   model.Id,
				Type: typeFromModel(model.Type),
			})
		}

		return result, nil
	case http.StatusInternalServerError:
		return nil, fmt.Errorf("server failure: %s", res.JSON500.Message)
	default:
		return nil, fmt.Errorf("unexpected response code: %d", res.StatusCode())
	}
}

func (storage *TypeStorage) Add(ctx context.Context, typ specification.TypeWithID) error {
	res, err := storage.httpClient.PostTasksSpecificationsTypesCreateWithResponse(ctx,
		oas.PostTasksSpecificationsTypesCreateJSONRequestBody{
			Id:   typ.ID,
			Type: typeToModel(typ.Type),
		})
	if err != nil {
		return fmt.Errorf("request failed: %w", err)
	}

	switch res.StatusCode() {
	case http.StatusOK:
		return nil
	case http.StatusInternalServerError:
		return fmt.Errorf("server failure: %s", res.JSON500.Message)
	default:
		return fmt.Errorf("unexpected response code: %d", res.StatusCode())
	}
}

func (storage *TypeStorage) Remove(ctx context.Context, id string) error {
	res, err := storage.httpClient.PostTasksSpecificationsTypesRemoveWithResponse(ctx,
		oas.PostTasksSpecificationsTypesRemoveJSONRequestBody{Id: id},
	)
	if err != nil {
		return fmt.Errorf("request failed: %w", err)
	}

	switch res.StatusCode() {
	case http.StatusOK:
		return nil
	case http.StatusInternalServerError:
		return fmt.Errorf("server failure: %s", res.JSON500.Message)
	default:
		return fmt.Errorf("unexpected response code: %d", res.StatusCode())
	}
}

func typeFromModel(model oas.Type) typing.Type {
	t := typing.Type{
		Name:     model.Name,
		SubTypes: map[string]typing.Type{},
	}

	for name, subType := range model.SubTypes.AdditionalProperties {
		t.SubTypes[name] = typeFromModel(subType)
	}

	return t
}

func resourceSetFromModel(model oas.SpecificationResourceSet) specification.ResourceSet {
	rs := specification.ResourceSet{
		Naming: map[int]string{},
		Types:  map[int]typing.Type{},
	}

	for _, obj := range model.Naming {
		rs.Naming[obj.Index] = obj.Name
	}

	for _, obj := range model.Types {
		rs.Types[obj.Index] = typeFromModel(obj.Type)
	}

	return rs
}

func specificationFromModel(model oas.Specification) specification.Specification {
	spec := specification.Specification{
		ID: model.Id,
		IO: specification.IO{
			Inputs:  resourceSetFromModel(model.Io.Inputs),
			Outputs: resourceSetFromModel(model.Io.Outputs),
		},
		Executable: specification.Executable{
			Type: model.Executable.Type,
			Data: model.Executable.Data,
		},
		Meta: model.Meta,
	}

	return spec
}

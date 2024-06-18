package http

import (
	"context"
	"fmt"
	"github.com/ischenkx/kantoku/pkg/common/data/record"
	"github.com/ischenkx/kantoku/pkg/common/data/record/ops"
	recutil "github.com/ischenkx/kantoku/pkg/common/data/record/util"
	"github.com/ischenkx/kantoku/pkg/core/resource"
	"github.com/ischenkx/kantoku/pkg/core/system"
	"github.com/ischenkx/kantoku/pkg/core/task"
	"github.com/ischenkx/kantoku/pkg/lib/gateway/api/http/oas"
	"github.com/ischenkx/kantoku/pkg/lib/tasks"
	"github.com/ischenkx/kantoku/pkg/lib/tasks/specification"
	"github.com/ischenkx/kantoku/pkg/lib/tasks/specification/typing"
	"github.com/samber/lo"
	"log/slog"
)

var _ oas.StrictServerInterface = (*Server)(nil)

type Server struct {
	system         system.AbstractSystem
	specifications *specification.Manager
}

func NewServer(system system.AbstractSystem, specifications *specification.Manager) *Server {
	return &Server{
		system:         system,
		specifications: specifications,
	}
}

func (server *Server) PostResourcesAllocate(ctx context.Context, request oas.PostResourcesAllocateRequestObject) (oas.PostResourcesAllocateResponseObject, error) {
	n := request.Params.Amount
	if n <= 0 {
		return oas.PostResourcesAllocate200JSONResponse{}, nil
	}

	ids, err := server.system.Resources().Alloc(ctx, n)
	if err != nil {
		return oas.PostResourcesAllocate500JSONResponse{Message: fmt.Sprintf("failed to allocate resources: %s", err)}, nil
	}

	return oas.PostResourcesAllocate200JSONResponse(ids), nil
}

func (server *Server) PostResourcesDeallocate(ctx context.Context, request oas.PostResourcesDeallocateRequestObject) (oas.PostResourcesDeallocateResponseObject, error) {
	err := server.system.Resources().Dealloc(ctx, *request.Body)
	if err != nil {
		return oas.PostResourcesDeallocate500JSONResponse{Message: fmt.Sprintf("failed to deallocate resources: %s", err)}, nil
	}

	return oas.PostResourcesDeallocate200JSONResponse{}, nil
}

func (server *Server) PostResourcesInitialize(ctx context.Context, request oas.PostResourcesInitializeRequestObject) (oas.PostResourcesInitializeResponseObject, error) {
	var resources []resource.Resource

	for _, initializer := range *request.Body {
		resources = append(resources, resource.Resource{
			Data: []byte(initializer.Value),
			ID:   initializer.Id,
		})
	}

	err := server.system.Resources().Init(ctx, resources)
	if err != nil {
		return oas.PostResourcesInitialize500JSONResponse{
			Message: fmt.Sprintf("failed to initialize resources: %s", err),
		}, nil
	}

	return oas.PostResourcesInitialize200JSONResponse{}, nil
}

func (server *Server) PostResourcesLoad(ctx context.Context, request oas.PostResourcesLoadRequestObject) (oas.PostResourcesLoadResponseObject, error) {
	resources, err := server.system.Resources().Load(ctx, *request.Body...)
	if err != nil {
		return oas.PostResourcesLoad500JSONResponse{
			Message: fmt.Sprintf("failed to load resources: %s", err),
		}, nil
	}

	return oas.PostResourcesLoad200JSONResponse(
		lo.Map(resources, func(res resource.Resource, _ int) oas.Resource {
			return oas.Resource{
				Id:     res.ID,
				Status: string(res.Status),
				Value:  string(res.Data),
			}
		})), nil
}

func (server *Server) PostTasksLoad(ctx context.Context, request oas.PostTasksLoadRequestObject) (oas.PostTasksLoadResponseObject, error) {
	tasks, err := recutil.List(
		ctx,
		server.system.
			Tasks().
			Filter(record.R{"id": ops.In[string](*request.Body...)}).
			Cursor().
			Iter(),
	)
	if err != nil {
		return oas.PostTasksLoad500JSONResponse{
			Message: fmt.Sprintf("failed to load tasks: %s", err),
		}, nil
	}

	dtoTasks := lo.Map(tasks, func(t task.Task, _ int) oas.Task {
		return TaskToDto(t)
	})

	return oas.PostTasksLoad200JSONResponse(dtoTasks), nil
}

func (server *Server) PostTasksSpawn(ctx context.Context, request oas.PostTasksSpawnRequestObject) (oas.PostTasksSpawnResponseObject, error) {
	spawnedTask, err := server.system.Spawn(ctx, task.Task{
		Inputs:  request.Body.Inputs,
		Outputs: request.Body.Outputs,
		Info:    request.Body.Info,
	})
	if err != nil {
		return oas.PostTasksSpawn500JSONResponse{
			Message: fmt.Sprintf("failed to spawn a new task: %s", err),
		}, nil
	}

	return oas.PostTasksSpawn200JSONResponse{Id: spawnedTask.ID}, nil
}

func (server *Server) PostTasksSpawnFromSpec(ctx context.Context, request oas.PostTasksSpawnFromSpecRequestObject) (oas.PostTasksSpawnFromSpecResponseObject, error) {

	type txParams struct {
		Inputs  []string
		Outputs []string
		Id      string
	}

	spec, err := server.specifications.Specifications().Get(ctx, request.Body.Specification)
	if err != nil {
		return oas.PostTasksSpawnFromSpec500JSONResponse{
			Message: fmt.Sprintf("failed to load specification: %s", err),
		}, nil
	}

	if len(request.Body.Parameters) != len(spec.IO.Inputs.Types) {
		return oas.PostTasksSpawnFromSpec500JSONResponse{
			Message: fmt.Sprintf("incorrect amount of input parameters: %d (expected %d)", len(request.Body.Parameters), len(spec.IO.Inputs.Types)),
		}, nil
	}

	tx := lo.NewTransaction[txParams]().
		Then(
			func(params txParams) (txParams, error) {
				// resource allocation
				inputAmounts := len(spec.IO.Inputs.Types)
				outputAmounts := len(spec.IO.Outputs.Types)
				totalResources := inputAmounts + outputAmounts
				resources, err := server.system.Resources().Alloc(ctx, totalResources)
				if err != nil {
					return params, fmt.Errorf("failed to allocate resources: %s", err)
				}
				params.Inputs = resources[:inputAmounts]
				params.Outputs = resources[inputAmounts:]

				return params, nil
			},
			func(params txParams) txParams {
				err := server.system.Resources().Dealloc(ctx, append(params.Inputs, params.Outputs...))
				if err != nil {
					slog.Error("failed to dealloc resources", slog.String("error", err.Error()))
				}
				return params
			},
		).
		Then(
			func(params txParams) (txParams, error) {
				initializedResources := make([]resource.Resource, 0, len(params.Inputs))
				for index := 0; index < len(params.Inputs); index++ {
					param := request.Body.Parameters[index]
					resourceId := params.Inputs[index]
					initializedResources = append(initializedResources, resource.Resource{
						Data: []byte(param),
						ID:   resourceId,
					})
				}

				err := server.system.Resources().Init(ctx, initializedResources)
				if err != nil {
					return params, fmt.Errorf("failed to initialize resources: %s", err)
				}

				return params, nil
			},
			func(params txParams) txParams {
				return params
			},
		).
		Then(
			func(params txParams) (txParams, error) {
				t, err := server.system.Spawn(ctx, task.New(
					task.WithInputs(params.Inputs...),
					task.WithOutputs(params.Outputs...),
					task.WithInfo(request.Body.Info),
					task.WithProperty("type", request.Body.Specification),
					tasks.DependOnInputs(),
				))
				if err != nil {
					return params, fmt.Errorf("failed to spawn a new task: %s", err)
				}

				params.Id = t.ID

				return params, nil
			},
			func(params txParams) txParams {
				return params
			},
		)

	result, err := tx.Process(txParams{})

	if err != nil {
		return oas.PostTasksSpawnFromSpec500JSONResponse{
			Message: fmt.Sprintf("failed to spawn a new task: %s", err),
		}, nil
	}

	return oas.PostTasksSpawnFromSpec200JSONResponse{Id: result.Id}, nil
}

func (server *Server) PostTasksCount(ctx context.Context, request oas.PostTasksCountRequestObject) (oas.PostTasksCountResponseObject, error) {
	records := server.system.Tasks()

	if request.Body.Filter != nil {
		filter := instantiateFilterOperators(*request.Body.Filter)
		records = records.Filter(filter.(map[string]any))
	}

	var cursor record.Cursor[task.Task]

	if request.Body.Cursor != nil {
		cursorConfig := request.Body.Cursor

		if cursorConfig.Distinct != nil {
			cursor = records.Distinct(*cursorConfig.Distinct...)
		} else {
			cursor = records.Cursor()
		}

		if cursorConfig.Sort != nil {
			var sorters []record.Sorter
			for _, val := range *cursorConfig.Sort {
				sorters = append(sorters, record.Sorter{
					Key:      val.Key,
					Ordering: record.Ordering(val.Ordering),
				})
			}
			cursor = cursor.Sort(sorters...)
		}

		if cursorConfig.Masks != nil {
			var masks []record.Mask
			for _, val := range *cursorConfig.Masks {
				masks = append(masks, record.Mask{
					Operation:       val.Operation,
					PropertyPattern: val.PropertyPattern,
				})
			}
			cursor = cursor.Mask(masks...)
		}

		if cursorConfig.Skip != nil {
			cursor = cursor.Skip(*cursorConfig.Skip)
		}

		if cursorConfig.Limit != nil {
			cursor = cursor.Limit(*cursorConfig.Limit)
		}
	} else {
		cursor = records.Cursor()
	}

	count, err := cursor.Count(ctx)
	if err != nil {
		return oas.PostTasksCount500JSONResponse{
			Message: fmt.Sprintf("failed to count: %s", err),
		}, nil
	}

	return oas.PostTasksCount200JSONResponse(count), nil
}

//func (server *Server) PostTasksInfoErase(ctx context.Context, request oas.PostTasksInfoEraseRequestObject) (oas.PostTasksInfoEraseResponseObject, error) {
//	var records record.Set = server.system.Info()
//
//	if request.Body.Filter != nil {
//		filter := instantiateFilterOperators(*request.Body.Filter)
//		records = records.Filter(filter.(map[string]any))
//	}
//
//	if err := records.Erase(ctx); err != nil {
//		return oas.PostTasksInfoErase500JSONResponse{
//			Message: fmt.Sprintf("failed to erase: %s", err),
//		}, nil
//	}
//
//	return oas.PostTasksInfoErase200JSONResponse{}, nil
//}
//
//func (server *Server) PostTasksInfoInsert(ctx context.Context, request oas.PostTasksInfoInsertRequestObject) (oas.PostTasksInfoInsertResponseObject, error) {
//	if err := server.system.Info().Insert(ctx, *request.Body); err != nil {
//		return oas.PostTasksInfoInsert500JSONResponse{
//			Message: fmt.Sprintf("failed to insert: %s", err),
//		}, nil
//	}
//
//	return oas.PostTasksInfoInsert200JSONResponse{}, nil
//}

func (server *Server) PostTasksFilter(ctx context.Context, request oas.PostTasksFilterRequestObject) (oas.PostTasksFilterResponseObject, error) {
	records := server.system.Tasks()

	if request.Body.Filter != nil {
		filter := instantiateFilterOperators(*request.Body.Filter)
		records = records.Filter(filter.(map[string]any))
	}

	var cursor record.Cursor[task.Task]

	if request.Body.Cursor != nil {
		cursorConfig := request.Body.Cursor

		if cursorConfig.Distinct != nil {
			cursor = records.Distinct(*cursorConfig.Distinct...)
		} else {
			cursor = records.Cursor()
		}

		if cursorConfig.Sort != nil {
			var sorters []record.Sorter
			for _, val := range *cursorConfig.Sort {
				sorters = append(sorters, record.Sorter{
					Key:      val.Key,
					Ordering: record.Ordering(val.Ordering),
				})
			}
			cursor = cursor.Sort(sorters...)
		}

		if cursorConfig.Masks != nil {
			var masks []record.Mask
			for _, val := range *cursorConfig.Masks {
				masks = append(masks, record.Mask{
					Operation:       val.Operation,
					PropertyPattern: val.PropertyPattern,
				})
			}
			cursor = cursor.Mask(masks...)
		}

		if cursorConfig.Skip != nil {
			cursor = cursor.Skip(*cursorConfig.Skip)
		}

		if cursorConfig.Limit != nil {
			cursor = cursor.Limit(*cursorConfig.Limit)
		}
	} else {
		cursor = records.Cursor()
	}

	tasks, err := recutil.List(ctx, cursor.Iter())
	if err != nil {
		return nil, fmt.Errorf("failed to load: %w", err)
	}

	dtos := lo.Map(tasks, func(t task.Task, _ int) oas.Task {
		return TaskToDto(t)
	})

	return oas.PostTasksFilter200JSONResponse(dtos), nil
}

func (server *Server) PostTasksUpdate(ctx context.Context, request oas.PostTasksUpdateRequestObject) (oas.PostTasksUpdateResponseObject, error) {
	records := server.system.Tasks()

	if request.Body.Filter != nil {
		filter := instantiateFilterOperators(request.Body.Filter)
		records = records.Filter(filter.(map[string]any))
	}

	var upsert record.R = nil
	if request.Body.Upsert != nil {
		upsert = *request.Body.Upsert
	}

	err := records.Update(ctx, request.Body.Update, upsert)
	if err != nil {
		return oas.PostTasksUpdate500JSONResponse{
			Message: fmt.Sprintf("failed to update: %s", err),
		}, nil
	}

	return oas.PostTasksUpdate200JSONResponse{}, nil
}

func instantiateFilterOperators(value any) any {
	switch d := value.(type) {
	case map[string]any:
		data, dataExists := d["Data"]
		typ, typeExists := d["Type"]
		if dataExists && typeExists {
			value = ops.Operation{
				Type: typ.(string),
				Data: instantiateFilterOperators(data),
			}
		} else {
			for key, v := range d {
				d[key] = instantiateFilterOperators(v)
			}
		}
	case []any:
		for index, v := range d {
			d[index] = instantiateFilterOperators(v)
		}
	}

	return value
}

func (server *Server) PostTasksSpecificationsCreate(ctx context.Context, request oas.PostTasksSpecificationsCreateRequestObject) (oas.PostTasksSpecificationsCreateResponseObject, error) {
	io := specification.IO{}

	// inputs
	{
		io.Inputs.Naming = map[int]string{}
		for _, obj := range request.Body.Io.Inputs.Naming {
			io.Inputs.Naming[obj.Index] = obj.Name
		}

		io.Inputs.Types = map[int]typing.Type{}
		for _, obj := range request.Body.Io.Inputs.Types {
			io.Inputs.Types[obj.Index] = parseType(obj.Type)
		}
	}

	// outputs
	{
		io.Outputs.Naming = map[int]string{}
		for _, obj := range request.Body.Io.Outputs.Naming {
			io.Outputs.Naming[obj.Index] = obj.Name
		}

		io.Outputs.Types = map[int]typing.Type{}
		for _, obj := range request.Body.Io.Outputs.Types {
			io.Outputs.Types[obj.Index] = parseType(obj.Type)
		}
	}

	spec := specification.Specification{
		ID: request.Body.Id,
		IO: io,
		Executable: specification.Executable{
			Type: request.Body.Executable.Type,
			Data: request.Body.Executable.Data,
		},
		Meta: request.Body.Meta,
	}
	if err := server.specifications.Specifications().Add(ctx, spec); err != nil {
		return oas.PostTasksSpecificationsCreate500JSONResponse{
			Message: err.Error(),
		}, nil
	}

	return oas.PostTasksSpecificationsCreate200Response{}, nil
}

func (server *Server) PostTasksSpecificationsGet(ctx context.Context, request oas.PostTasksSpecificationsGetRequestObject) (oas.PostTasksSpecificationsGetResponseObject, error) {
	spec, err := server.specifications.Specifications().Get(ctx, request.Body.Id)
	if err != nil {
		return oas.PostTasksSpecificationsGet500JSONResponse{
			Message: err.Error(),
		}, nil
	}

	return oas.PostTasksSpecificationsGet200JSONResponse(specificationToModel(spec)), nil
}

func (server *Server) PostTasksSpecificationsGetAll(ctx context.Context, request oas.PostTasksSpecificationsGetAllRequestObject) (oas.PostTasksSpecificationsGetAllResponseObject, error) {
	specs, err := server.specifications.Specifications().GetAll(ctx)
	if err != nil {
		return oas.PostTasksSpecificationsGetAll500JSONResponse{
			Message: err.Error(),
		}, nil
	}

	specModels := make([]oas.Specification, 0, len(specs))
	for _, spec := range specs {
		specModels = append(specModels, specificationToModel(spec))
	}

	return oas.PostTasksSpecificationsGetAll200JSONResponse(specModels), nil
}

func (server *Server) PostTasksSpecificationsRemove(ctx context.Context, request oas.PostTasksSpecificationsRemoveRequestObject) (oas.PostTasksSpecificationsRemoveResponseObject, error) {
	if err := server.specifications.Specifications().Remove(ctx, request.Body.Id); err != nil {
		return oas.PostTasksSpecificationsRemove500JSONResponse{
			Message: err.Error(),
		}, nil
	}

	return oas.PostTasksSpecificationsRemove200Response{}, nil
}

func (server *Server) PostTasksSpecificationsTypesCreate(ctx context.Context, request oas.PostTasksSpecificationsTypesCreateRequestObject) (oas.PostTasksSpecificationsTypesCreateResponseObject, error) {
	typ := specification.TypeWithID{
		ID:   request.Body.Id,
		Type: parseType(request.Body.Type),
	}

	if err := server.specifications.Types().Add(ctx, typ); err != nil {
		return oas.PostTasksSpecificationsTypesCreate500JSONResponse{
			Message: err.Error(),
		}, nil
	}

	return oas.PostTasksSpecificationsTypesCreate200Response{}, nil
}

func (server *Server) PostTasksSpecificationsTypesGet(ctx context.Context, request oas.PostTasksSpecificationsTypesGetRequestObject) (oas.PostTasksSpecificationsTypesGetResponseObject, error) {
	typ, err := server.specifications.Types().Get(ctx, request.Body.Id)
	if err != nil {
		return oas.PostTasksSpecificationsTypesGet500JSONResponse{
			Message: err.Error(),
		}, nil
	}

	return oas.PostTasksSpecificationsTypesGet200JSONResponse{
		Id:   typ.ID,
		Type: typeToModel(typ.Type),
	}, nil
}

func (server *Server) PostTasksSpecificationsTypesGetAll(ctx context.Context, request oas.PostTasksSpecificationsTypesGetAllRequestObject) (oas.PostTasksSpecificationsTypesGetAllResponseObject, error) {
	types, err := server.specifications.Types().GetAll(ctx)
	if err != nil {
		return oas.PostTasksSpecificationsTypesGetAll500JSONResponse{
			Message: err.Error(),
		}, nil
	}

	typeModels := make([]oas.TypeWithID, 0, len(types))
	for _, typ := range types {
		typeModels = append(typeModels, oas.TypeWithID{
			Id:   typ.ID,
			Type: typeToModel(typ.Type),
		})
	}

	return oas.PostTasksSpecificationsTypesGetAll200JSONResponse(typeModels), nil
}

func (server *Server) PostTasksSpecificationsTypesRemove(ctx context.Context, request oas.PostTasksSpecificationsTypesRemoveRequestObject) (oas.PostTasksSpecificationsTypesRemoveResponseObject, error) {
	if err := server.specifications.Types().Remove(ctx, request.Body.Id); err != nil {
		return oas.PostTasksSpecificationsTypesRemove500JSONResponse{
			Message: err.Error(),
		}, nil
	}

	return oas.PostTasksSpecificationsTypesRemove200Response{}, nil
}

func parseType(t oas.Type) typing.Type {
	result := typing.Type{
		Name:     t.Name,
		SubTypes: map[string]typing.Type{},
	}

	for key, subType := range t.SubTypes.AdditionalProperties {
		result.SubTypes[key] = parseType(subType)
	}

	return result
}

func typeToModel(t typing.Type) oas.Type {
	result := oas.Type{
		Name: t.Name,
		SubTypes: oas.Type_SubTypes{
			AdditionalProperties: map[string]oas.Type{},
		},
	}

	for key, value := range t.SubTypes {
		result.SubTypes.AdditionalProperties[key] = typeToModel(value)
	}

	return result
}

func resourceSetToModel(rs specification.ResourceSet) oas.SpecificationResourceSet {
	var result oas.SpecificationResourceSet

	type namingT struct {
		Index int    `json:"index"`
		Name  string `json:"name"`
	}

	type typeT struct {
		Index int      `json:"index"`
		Type  oas.Type `json:"type"`
	}

	for index, name := range rs.Naming {
		result.Naming = append(result.Naming, namingT{
			Index: index,
			Name:  name,
		})
	}

	for index, typ := range rs.Types {
		result.Types = append(result.Types, typeT{
			Index: index,
			Type:  typeToModel(typ),
		})
	}

	return result
}

func specificationToModel(spec specification.Specification) oas.Specification {
	return oas.Specification{
		Executable: oas.SpecificationExecutable{
			Data: spec.Executable.Data,
			Type: spec.Executable.Type,
		},
		Id: spec.ID,
		Io: oas.SpecificationIO{
			Inputs:  resourceSetToModel(spec.IO.Inputs),
			Outputs: resourceSetToModel(spec.IO.Outputs),
		},
		Meta: spec.Meta,
	}
}

package server

import (
	"context"
	"fmt"
	"github.com/ischenkx/kantoku/pkg/common/data/record"
	"github.com/ischenkx/kantoku/pkg/common/data/record/ops"
	recutil "github.com/ischenkx/kantoku/pkg/common/data/record/util"
	"github.com/ischenkx/kantoku/pkg/core/resource"
	"github.com/ischenkx/kantoku/pkg/core/system"
	"github.com/ischenkx/kantoku/pkg/core/task"
	"github.com/ischenkx/kantoku/pkg/lib/gateway/api/http/converters"
	"github.com/ischenkx/kantoku/pkg/lib/gateway/api/http/oas"
	"github.com/samber/lo"
)

var _ oas.StrictServerInterface = (*Server)(nil)

type Server struct {
	system system.AbstractSystem
}

func New(system system.AbstractSystem) *Server {
	return &Server{system: system}
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
		return converters.TaskToDto(t)
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
		return converters.TaskToDto(t)
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

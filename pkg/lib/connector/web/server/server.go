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
	"github.com/labstack/echo/v4"
	"github.com/samber/lo"
	"net/http"
	"sort"
	"strconv"
	"strings"
	"time"
)

type Server struct {
	System system.AbstractSystem
}

func (server *Server) Echo() *echo.Echo {
	// Create a new Echo instance
	e := echo.New()

	// Create groups for "users" and "groups"
	resourcesGroup := e.Group("/resources")
	tasksGroup := e.Group("/tasks")

	// Routes for "users" group
	resourcesGroup.GET("", server.getResources)
	resourcesGroup.GET("/:id", server.getResource)
	resourcesGroup.POST("", server.allocateResources)
	resourcesGroup.PATCH("/:id", server.initializeResource)
	resourcesGroup.DELETE("/:id", server.deleteResource)

	// Routes for "groups" group
	tasksGroup.GET("", server.getTasks)
	tasksGroup.GET("/:id", server.getTask)
	tasksGroup.POST("", server.createTask)
	tasksGroup.PATCH("/:id", server.updateTask)
	tasksGroup.DELETE("/:id", server.deleteTask)

	return e
}

func (server *Server) getResources(c echo.Context) error {
	ctx := c.Request().Context()

	var params ListParams[struct{}]
	if err := c.Bind(&params); err != nil {
		return c.JSON(http.StatusBadRequest, fmt.Sprintf("Error parsing query parameters: %s", err))
	}

	if len(params.IDs) == 0 {
		return c.JSON(http.StatusOK, []ResourceDto{})
	}

	resources, err := server.System.Resources().Load(ctx, params.IDs...)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, fmt.Sprintf("Error loading resources: %s", err))
	}

	if params.Sort != nil && params.Order != nil {
		keys := strings.Split(*params.Sort, ",")
		orders := strings.Split(*params.Order, ",")

		if len(keys) != len(orders) {
			return c.JSON(http.StatusBadRequest, fmt.Sprintf("Keys amount does not match amount of orders"))
		}

		sorters := lo.Zip2(keys, orders)

		sort.Slice(resources, func(i, j int) bool {
			res1, res2 := &resources[i], &resources[j]

			for _, sorter := range sorters {
				key, order := sorter.A, sorter.B

				shouldInvert := order == "desc"

				var value1, value2 string
				switch key {
				case "id":
					value1, value2 = res1.ID, res2.ID
				case "status":
					value1, value2 = string(res1.Status), string(res2.Status)
				}

				less := value1 < value2
				if shouldInvert {
					less = !less
				}

				if less {
					return true
				}
			}

			return false
		})
	}

	startIndex, endIndex := 0, len(resources)

	if params.Start != nil {
		startIndex = max(startIndex, *params.Start)
	}

	if params.End != nil {
		endIndex = min(endIndex, *params.End)
	}

	resources = resources[startIndex:endIndex]

	dtos := lo.Map(resources, func(res resource.Resource, _ int) ResourceDto {
		return ResourceDto{
			ID:     res.ID,
			Status: string(res.Status),
			Data:   string(res.Data),
		}
	})

	return c.JSON(http.StatusOK, dtos)
}

func (server *Server) getResource(c echo.Context) error {
	ctx := c.Request().Context()

	id := c.Param("id")

	resources, err := server.System.Resources().Load(ctx, id)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, fmt.Sprintf("Error loading resources: %s", err))
	}

	if len(resources) != 1 {
		return c.JSON(http.StatusInternalServerError, fmt.Sprintf("Unexpected amount of loaded resources"))
	}

	res := resources[0]

	return c.JSON(http.StatusOK, ResourceDto{
		ID:     res.ID,
		Status: string(res.Status),
		Data:   string(res.Data),
	})
}

func (server *Server) allocateResources(c echo.Context) error {
	ctx := c.Request().Context()

	var request AllocateResourceRequest
	if err := c.Bind(&request); err != nil {
		return c.JSON(http.StatusBadRequest, fmt.Sprintf("Error parsing query parameters: %s", err))
	}

	resources, err := server.System.Resources().Alloc(ctx, request.Amount)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, fmt.Sprintf("Error allocating resources: %s", err))
	}

	dtos := lo.Map(resources, func(id string, _ int) ResourceDto {
		return ResourceDto{
			ID:     id,
			Status: resource.Allocated,
			Data:   "",
		}
	})

	return c.JSON(http.StatusOK, dtos)
}

func (server *Server) initializeResource(c echo.Context) error {
	ctx := c.Request().Context()

	var request InitializeResourceRequest
	if err := c.Bind(&request); err != nil {
		return c.JSON(http.StatusBadRequest, fmt.Sprintf("Error parsing query and body: %s", err))
	}

	res := resource.Resource{
		Data: []byte(request.Data),
		ID:   request.ID,
	}

	if err := server.System.Resources().Init(ctx, []resource.Resource{res}); err != nil {
		return c.JSON(http.StatusInternalServerError, fmt.Sprintf("Error initializing the resource: %s", err))
	}

	return c.JSON(http.StatusOK, nil)
}

func (server *Server) deleteResource(c echo.Context) error {
	return c.JSON(http.StatusMethodNotAllowed, nil)
}

func (server *Server) getTasks(c echo.Context) error {
	//ctx := c.Request().Context()

	ctx := context.Background()

	type other struct {
		Status       []string   `query:"info.status"`
		Context      []string   `query:"context"`
		UpdatedAtGte *time.Time `query:"info.updatedAt_gte"`
		UpdatedAtLte *time.Time `query:"info.updatedAt_lte"`
	}

	var params ListParams[other]
	if err := c.Bind(&params); err != nil {
		fmt.Println(err)
		return c.JSON(http.StatusBadRequest, fmt.Sprintf("Error parsing query parameters: %s", err))
	}

	set := server.
		System.
		Tasks().
		Filter(record.R{"id": ops.Not(ops.In[any](nil))})

	if len(params.IDs) > 0 {
		set = set.Filter(record.R{
			"id": ops.In[string](params.IDs...),
		})
	}

	if len(params.Other.Status) > 0 {
		set = set.Filter(record.R{
			"status": ops.In[string](params.Other.Status...),
		})
	}

	if len(params.Other.Context) > 0 {
		set = set.Filter(record.R{
			"context_id": ops.In[string](params.Other.Context...),
		})
	}

	if params.Other.UpdatedAtLte != nil || params.Other.UpdatedAtGte != nil {
		var filters []any

		if params.Other.UpdatedAtGte != nil {
			filters = append(filters, ops.GreaterOrEqualThan(params.Other.UpdatedAtGte.Unix()))
		}

		if params.Other.UpdatedAtLte != nil {
			filters = append(filters, ops.LessOrEqualThan(params.Other.UpdatedAtLte.Unix()))
		}

		set = set.Filter(record.R{
			"updated_at": ops.And(filters...),
		})
	}

	cursor := set.Cursor()

	if params.Sort != nil && params.Order != nil {
		keys := strings.Split(*params.Sort, ",")
		orders := strings.Split(*params.Order, ",")

		if len(keys) != len(orders) {
			return c.JSON(http.StatusBadRequest, fmt.Sprintf("Keys amount does not match amount of orders"))
		}

		rawSorters := lo.Zip2(keys, orders)

		sorters := lo.Map(rawSorters, func(pair lo.Tuple2[string, string], _ int) (sorter record.Sorter) {
			switch pair.B {
			case "desc":
				sorter.Ordering = record.Desc
			case "asc":
				sorter.Ordering = record.Asc
			default:
				sorter.Ordering = record.Asc
			}

			sorter.Key = pair.A

			return
		})

		cursor = cursor.Sort(sorters...)
	} else {
		cursor = cursor.Sort(
			record.Sorter{
				Key:      "id",
				Ordering: record.Asc,
			},
		)
	}

	totalCount, err := cursor.Count(ctx)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, fmt.Sprintf("Failed to count tasks: %s", err))
	}

	if params.Start != nil {
		cursor = cursor.Skip(*params.Start)
	}

	if params.End != nil {
		start := 0
		if params.Start != nil {
			start = *params.Start
		}

		limit := *params.End - start

		cursor = cursor.Limit(limit)
	}

	tasks, err := recutil.List(ctx, cursor.Iter())
	if err != nil {
		return fmt.Errorf("failed to load data: %w", err)
	}

	dtos := lo.FilterMap(tasks, func(t task.Task, _ int) (TaskDto, bool) {
		return TaskDto{
			ID:      t.ID,
			Inputs:  t.Inputs,
			Outputs: t.Outputs,
			Info:    t.Info,
		}, true
	})

	c.Response().Header().Set("X-Total-Count", strconv.Itoa(totalCount))
	c.Response().Header().Set("Access-Control-Expose-Headers", "X-Total-Count")

	return c.JSON(http.StatusOK, dtos)
}

func (server *Server) getTask(c echo.Context) error {
	ctx := context.Background()

	id := c.Param("id")

	t, err := recutil.Single(ctx, server.System.
		Tasks().
		Filter(record.R{"id": id}).
		Cursor().
		Iter(),
	)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, fmt.Sprintf("Failed to parse task from record: %s", err))
	}

	dto := TaskDto{
		ID:      t.ID,
		Inputs:  t.Inputs,
		Outputs: t.Outputs,
		Info:    t.Info,
	}

	return c.JSON(http.StatusOK, dto)

}

func (server *Server) createTask(c echo.Context) error {
	return c.JSON(http.StatusCreated, "Task created")
}

func (server *Server) updateTask(c echo.Context) error {
	id := c.Param("id")
	return c.JSON(http.StatusOK, fmt.Sprintf("Task updated with ID %s", id))
}

func (server *Server) deleteTask(c echo.Context) error {
	id := c.Param("id")
	return c.JSON(http.StatusOK, fmt.Sprintf("Task deleted with ID %s", id))
}

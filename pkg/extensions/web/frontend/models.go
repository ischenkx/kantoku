package frontend

import (
	"github.com/ischenkx/kantoku/pkg/common/data/record"
	"github.com/ischenkx/kantoku/pkg/system/kernel/task"
)

type ListParams[Other any] struct {
	IDs   []string `query:"id"`
	Start *int     `query:"_start"`
	End   *int     `query:"_end"`
	Sort  *string  `query:"_sort"`
	Order *string  `query:"_order"`
	Other Other
}

type ResourceDto struct {
	ID     string `json:"id"`
	Status string `json:"status"`
	Data   string `json:"data"`
}

type TaskDto struct {
	ID         string          `json:"id"`
	Inputs     []string        `json:"inputs"`
	Outputs    []string        `json:"outputs"`
	Properties task.Properties `json:"properties"`
	Info       record.R        `json:"info"`
}

type AllocateResourceRequest struct {
	Amount int `json:"amount"`
}

type InitializeResourceRequest struct {
	ID   string `param:"id"`
	Data string `json:"data"`
}

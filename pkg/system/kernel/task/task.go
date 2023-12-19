package task

import (
	"github.com/ischenkx/kantoku/pkg/system/kernel/resource"
)

type Properties struct {
	Data map[string]any
	Sub  map[string]Properties
}

type Task struct {
	Inputs     []resource.ID
	Outputs    []resource.ID
	Properties Properties
	ID         string
}

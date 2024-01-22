package task

import (
	"github.com/ischenkx/kantoku/pkg/common/data/record"
	"github.com/ischenkx/kantoku/pkg/core/resource"
)

type Task struct {
	Inputs  []resource.ID
	Outputs []resource.ID
	ID      string
	Info    record.R
}

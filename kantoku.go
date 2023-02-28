package kantoku

import (
	"kantoku/common/db/kv"
	"kantoku/common/deps"
	"kantoku/core/l0/cell"
	"kantoku/core/l0/event"
)

type Kantoku struct {
	events event.Bus
	deps   deps.DB
	tasks  kv.Database[Task]
	cells  cell.Storage
}

func New(events event.Bus, deps deps.DB, tasks kv.Database[Task], cells cell.Storage) *Kantoku {
	return &Kantoku{
		events: events,
		deps:   deps,
		tasks:  tasks,
		cells:  cells,
	}
}

func (kantoku *Kantoku) Events() event.Bus {
	return kantoku.events
}

func (kantoku *Kantoku) Deps() deps.DB {
	return kantoku.deps
}

func (kantoku *Kantoku) Tasks() kv.Database[Task] {
	return kantoku.tasks
}

func (kantoku *Kantoku) Cells() cell.Storage {
	return kantoku.cells
}

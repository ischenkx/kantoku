package storage

import "context"

type Settings struct {
	Type string
	Meta map[string]any
}

type Param struct {
	Name  string
	Value any
}

type Command struct {
	Operation string
	Params    []Param
	Meta      map[string]any
}

type Document = map[string]any

type Storage interface {
	Settings(ctx context.Context) (Settings, error)
	Exec(ctx context.Context, command Command) ([]Document, error)
}

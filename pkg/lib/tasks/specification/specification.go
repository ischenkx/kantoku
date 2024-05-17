package specification

import "github.com/ischenkx/kantoku/pkg/lib/tasks/specification/typing"

type ResourceSet struct {
	Naming map[int]string
	Types  map[int]typing.Type
}

type IO struct {
	Inputs  ResourceSet
	Outputs ResourceSet
}

type Executable struct {
	Type string
	Data map[string]any
}

type Specification struct {
	// ID is a path to specification, it's expected to be like /a/b/c/task.go
	ID         string
	IO         IO
	Executable Executable
	Meta       map[string]any
}

type TypeWithID struct {
	ID   string
	Type typing.Type
}

package resource

type Status string

const (
	DoesNotExist Status = "does_not_exist"
	Allocated           = "allocated"
	Ready               = "ready"
)

type ID = string

type Resource struct {
	Data   []byte
	ID     ID
	Status Status
}

package deps

type Status string

const (
	OK           Status = "ok"
	Failed              = "failed"
	Pending             = "pending"
	DoesNotExist        = "does_not_exist"
)

type Dependency struct {
	ID     string
	Status Status
}

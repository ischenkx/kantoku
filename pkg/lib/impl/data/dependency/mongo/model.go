package mongorep

const (
	DependencyKind = "dependency"
	GroupKind      = "group"
)

const (
	GroupCreated      = "created"
	GroupInitializing = "initializing"
	GroupInitialized  = "initialized"
)

type Doc struct {
	ID        string `bson:"_id"`
	ContextID string `bson:"context_id"`
	Kind      string `bson:"kind"`
	// UnixNano time
	UpdatedAt int64 `bson:"updated_at"`
}

type Dependency struct {
	Doc

	GroupsProcessed bool   `bson:"groups_processed"`
	Status          string `bson:"status"`
}

type Set[T comparable] map[T]*struct{}

func (s Set[T]) Add(values ...T) {
	for _, value := range values {
		s[value] = nil
	}
}

func (s Set[T]) Remove(values ...T) {
	for _, value := range values {
		delete(s, value)
	}
}

type Group struct {
	Doc
	Status  string      `bson:"status"`
	Pending Set[string] `bson:"pending"`
	Ready   Set[string] `bson:"ready"`
}

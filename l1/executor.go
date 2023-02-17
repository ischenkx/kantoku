package l1

type Executor interface {
	Execute(Task) (Result, error)
}

package kantoku

type Initializeable interface {
	Initialize(ctx *Context) (any, error)
}

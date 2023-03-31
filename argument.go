package kantoku

type ArgumentInitializer interface {
	Initialize(ctx *Context) (any, error)
}

package kantoku

type Plugin interface {
	Initialize(kantoku *Kantoku)
	BeforeInitialized(ctx *Context) error
	AfterInitialized(ctx *Context)
	BeforeScheduled(ctx *Context) error
	AfterScheduled(ctx *Context)
}

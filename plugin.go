package kantoku

type Plugin any

type InitializePlugin interface {
	Initialize(kantoku Kantoku)
}

type BeforeInitializedPlugin interface {
	BeforeInitialized(ctx *Context) error
}

type AfterInitializedPlugin interface {
	AfterInitialized(ctx *Context)
}

type BeforeScheduledPlugin interface {
	BeforeScheduled(ctx *Context) error
}

type AfterScheduledPlugin interface {
	AfterScheduled(ctx *Context)
}

// initializer, ok := plugin.(InitializePlugin); ok { ... }

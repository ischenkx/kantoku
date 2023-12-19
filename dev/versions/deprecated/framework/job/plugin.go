package job

type Plugin any

type InitializePlugin interface {
	Initialize(kernel *Manager)
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

type BeforePluginInitPlugin interface {
	BeforePluginInit(kernel *Manager, plugin Plugin) error
}

type AfterPluginInitPlugin interface {
	AfterPluginInit(kernel *Manager, plugin Plugin)
}

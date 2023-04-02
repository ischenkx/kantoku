package kantoku

type Scheduler interface {
	Schedule(ctx *Context) error
}

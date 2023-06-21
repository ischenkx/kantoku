package context

type Context struct {
	ID     string
	Parent string
}

var Empty = Context{}

package kantoku

import (
	"context"
	"time"
)

type PluginData struct {
	data map[string]any
}

func NewPluginData() PluginData {
	return PluginData{data: map[string]any{}}
}

func (pd PluginData) Get(key string) (any, bool) {
	val, ok := pd.data[key]
	return val, ok
}

func (pd PluginData) Set(key string, value any) {
	pd.data[key] = value
}

func (pd PluginData) Del(key string) {
	delete(pd.data, key)
}

func (pd PluginData) GetOrSet(key string, value func() any) any {
	if val, ok := pd.Get(key); ok {
		return val
	}
	val := value()
	pd.Set(key, val)
	return val
}

type Context struct {
	Kantoku *Kantoku
	Task    *TaskInstance
	Log     *Log
	data    PluginData
	parent  context.Context
}

func NewContext(ctx context.Context) *Context {
	return &Context{
		parent: ctx,
		Log:    &Log{},
		data:   NewPluginData(),
	}
}

func (c *Context) Data() PluginData {
	return c.data
}

func (c *Context) Deadline() (deadline time.Time, ok bool) {
	return c.parent.Deadline()
}

func (c *Context) Done() <-chan struct{} {
	return c.parent.Done()
}

func (c *Context) Err() error {
	return c.parent.Err()
}

func (c *Context) Value(key any) any {
	return c.parent.Value(key)
}

func (c *Context) finalize() {
	c.tryPropagate()
	c.Task = nil
	c.Kantoku = nil

	// might be harmful
	//c.parent = nil
}

func (c *Context) tryPropagate() {
	if parent, ok := c.parent.(*Context); ok && c.Kantoku != nil {
		parent.Log.Merge(c.Log)
	}
}

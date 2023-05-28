package taskdep

import "kantoku"

func Dep(task string) kantoku.Option {
	return func(ctx *kantoku.Context) error {
		data := ctx.Data().GetOrSet("taskdep", func() any { return &PluginData{} }).(*PluginData)
		data.Subtasks = append(data.Subtasks, task)

		return nil
	}
}

//func Output(task string) kantoku.Data {
//
//}
//
//type output struct {
//	id string
//}
//
//func (o output) Initialize(ctx *kantoku.Context) ([]byte, error) {
//	//TODO implement me
//	panic("implement me")
//}

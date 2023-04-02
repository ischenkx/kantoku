package kantoku

import "context"

type PropertyEvaluator interface {
	Evaluate(ctx context.Context, task string) (any, error)
}

type Properties struct {
	evaluator     PropertyEvaluator
	subProperties map[string]*Properties
}

func NewProperties() *Properties {
	return &Properties{
		subProperties: map[string]*Properties{},
	}
}

func (props *Properties) Get(path ...string) (PropertyEvaluator, bool) {
	if len(path) == 0 {
		return props.evaluator, props.evaluator == nil
	}

	subProps, ok := props.subProperties[path[0]]
	if !ok {
		return nil, false
	}

	return subProps.Get(path[1:]...)
}

func (props *Properties) Del(path ...string) {
	if len(path) == 0 {
		props.evaluator = nil
		return
	}

	if subProps, ok := props.subProperties[path[0]]; ok {
		subProps.Del(path[1:]...)
	}
}

func (props *Properties) Set(evaluator PropertyEvaluator, path ...string) {
	if len(path) == 0 {
		props.evaluator = evaluator
		return
	}

	subProps, ok := props.subProperties[path[0]]
	if !ok {
		subProps = NewProperties()
		props.subProperties[path[0]] = subProps
	}

	subProps.Set(evaluator, path[1:]...)
}

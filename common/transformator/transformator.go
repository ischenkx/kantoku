package transformator

type Transformator[In, Out any] func(In) Out

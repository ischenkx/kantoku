package transformer

type Transformer[In, Out any] func(In) Out

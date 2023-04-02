package transformer

type Transformer[In, Out any] func(In) Out
type CheckedTransformer[In, Out any] func(In) (Out, bool)

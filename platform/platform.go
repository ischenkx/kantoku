package platform

type Platform[T Task] struct {
	inputs  Inputs[T]
	outputs Outputs
	broker  Broker
}

func New[T Task](inputs Inputs[T], outputs Outputs, broker Broker) Platform[T] {
	return Platform[T]{
		inputs:  inputs,
		outputs: outputs,
		broker:  broker,
	}
}

func (p Platform[T]) Inputs() Inputs[T] {
	return p.inputs
}

func (p Platform[T]) Outputs() Outputs {
	return p.outputs
}

func (p Platform[T]) Broker() Broker {
	return p.broker
}

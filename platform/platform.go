package platform

type Platform[T Task] struct {
	db      DB[T]
	inputs  Inputs
	outputs Outputs
	broker  Broker
}

func New[T Task](db DB[T], inputs Inputs, outputs Outputs, broker Broker) Platform[T] {
	return Platform[T]{
		db:      db,
		inputs:  inputs,
		outputs: outputs,
		broker:  broker,
	}
}

func (p Platform[T]) DB() DB[T] {
	return p.db
}

func (p Platform[T]) Inputs() Inputs {
	return p.inputs
}

func (p Platform[T]) Outputs() Outputs {
	return p.outputs
}

func (p Platform[T]) Broker() Broker {
	return p.broker
}

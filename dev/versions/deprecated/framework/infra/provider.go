package infra

type Provider interface {
	Demons() []Demon
}

package evexec

import "errors"

type TopicResolver interface {
	Resolve(eventName string) (topic string, err error)
}

type MapResolver map[string]string

func (resolver MapResolver) Resolve(name string) (string, error) {
	topic, ok := resolver[name]
	if !ok {
		return "", errors.New("failed to find a value for the given name")
	}

	return topic, nil
}

type ConstantResolver string

func (resolver ConstantResolver) Resolve(string) (string, error) {
	return string(resolver), nil
}

package main

import (
	"github.com/ischenkx/kantoku/pkg/lib/tasks/functional"
	"github.com/ischenkx/kantoku/pkg/lib/tasks/future"
	"reflect"
)

type RequestTask struct {
}

type RequestInput struct {
	url future.Future[string]
}

type RequestOutput struct {
	body future.Future[string]
}

func (r RequestTask) EmptyOutput() RequestOutput {
	return RequestOutput{body: future.Empty[string]()}
}

func (r RequestTask) Call(context functional.Context, input RequestInput) (RequestOutput, error) {
	//todo
	return RequestOutput{}, nil
}

func (r RequestTask) InputType() reflect.Type {
	return reflect.TypeOf(RequestInput{})
}

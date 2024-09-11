package core

import (
	"context"
	"errors"
	"fmt"
	"github.com/ischenkx/kantoku/pkg/common/data/storage"
	"github.com/mitchellh/mapstructure"
)

type Task struct {
	Inputs  []string
	Outputs []string
	ID      string
	Info    map[string]any
}

type Option func(t *Task)

func New(options ...Option) Task {
	t := Task{Info: map[string]any{}}

	for _, option := range options {
		option(&t)
	}

	return t
}

func (task Task) ContextID() string {
	rawContextID, ok := task.Info["context_id"]
	if !ok {
		return ""
	}

	return rawContextID.(string)
}

func (task Task) AsDoc() storage.Document {
	return map[string]any{
		"id":      task.ID,
		"inputs":  task.Inputs,
		"outputs": task.Outputs,
		"info":    task.Info,
	}
}

type TaskCodec struct{}

func (c TaskCodec) Encode(task Task) (map[string]any, error) {
	return task.AsDoc(), nil
}

func (c TaskCodec) Decode(rec map[string]any) (Task, error) {
	if rec == nil {
		return Task{}, errors.New("record is nil")
	}
	var t Task

	if err := mapstructure.Decode(rec, &t); err != nil {
		return Task{}, fmt.Errorf("failed to decode: %w", err)
	}

	return t, nil
}

type TaskDB interface {
	storage.Storage

	Insert(ctx context.Context, tasks []Task) error
	Delete(ctx context.Context, ids []string) error
	ByIDs(ctx context.Context, ids []string) ([]Task, error)
	UpdateByIDs(ctx context.Context, ids []string, properties map[string]any) error
	GetWithProperties(ctx context.Context, propertiesToValues map[string][]any) ([]Task, error)
	UpdateWithProperties(ctx context.Context, propertiesToValues map[string][]any, newProperties map[string]any) (updatedDocs int, err error)
}

package resources

import (
	"context"
	"fmt"
	"github.com/ischenkx/kantoku/pkg/core/resource"
)

type Observer interface {
	BeforeLoad(ctx context.Context, ids []resource.ID) error
	AfterLoad(ctx context.Context, resources []resource.Resource)
	OnLoadError(ctx context.Context, ids []resource.ID, err error)

	BeforeAlloc(ctx context.Context, n int) error
	AfterAlloc(ctx context.Context, resources []resource.ID)
	OnAllocError(ctx context.Context, n int, err error)

	BeforeInit(ctx context.Context, resources []resource.Resource) error
	AfterInit(ctx context.Context, resources []resource.Resource)
	OnInitError(ctx context.Context, resources []resource.Resource, err error)

	BeforeDealloc(ctx context.Context, resources []resource.ID) error
	AfterDealloc(ctx context.Context, resources []resource.ID)
	OnDeallocError(ctx context.Context, resources []resource.ID, err error)
}

func Observe(raw resource.Storage, observer Observer) *Observable {
	return &Observable{
		raw:      raw,
		observer: observer,
	}
}

type Observable struct {
	raw      resource.Storage
	observer Observer
}

func (storage *Observable) Load(ctx context.Context, ids ...resource.ID) ([]resource.Resource, error) {
	if err := storage.observer.BeforeLoad(ctx, ids); err != nil {
		return nil, fmt.Errorf("observer.beforeLoad failed: %w", err)
	}

	resources, err := storage.raw.Load(ctx, ids...)
	if err != nil {
		storage.observer.OnLoadError(ctx, ids, err)
		return nil, err
	}

	storage.observer.AfterLoad(ctx, resources)

	return resources, nil
}

func (storage *Observable) Alloc(ctx context.Context, amount int) ([]resource.ID, error) {
	if err := storage.observer.BeforeAlloc(ctx, amount); err != nil {
		return nil, fmt.Errorf("observer.beforeAlloc failed: %w", err)
	}

	ids, err := storage.raw.Alloc(ctx, amount)
	if err != nil {
		storage.observer.OnAllocError(ctx, amount, err)
		return nil, err
	}

	storage.observer.AfterAlloc(ctx, ids)

	return ids, nil
}

func (storage *Observable) Init(ctx context.Context, resources []resource.Resource) error {
	if err := storage.observer.BeforeInit(ctx, resources); err != nil {
		return fmt.Errorf("observer.beforeInit failed: %w", err)
	}

	err := storage.raw.Init(ctx, resources)
	if err != nil {
		storage.observer.OnInitError(ctx, resources, err)
		return nil
	}

	storage.observer.AfterInit(ctx, resources)

	return nil
}

func (storage *Observable) Dealloc(ctx context.Context, ids []resource.ID) error {
	if err := storage.observer.BeforeDealloc(ctx, ids); err != nil {
		return fmt.Errorf("observer.beforeDealloc failed: %w", err)
	}

	err := storage.raw.Dealloc(ctx, ids)
	if err != nil {
		storage.observer.OnDeallocError(ctx, ids, err)
		return nil
	}

	storage.observer.AfterDealloc(ctx, ids)

	return nil
}

type FunctionalObserver struct {
	BeforeLoadF     func(ctx context.Context, ids []resource.ID) error
	AfterLoadF      func(ctx context.Context, resources []resource.Resource)
	OnLoadErrorF    func(ctx context.Context, ids []resource.ID, err error)
	BeforeAllocF    func(ctx context.Context, n int) error
	AfterAllocF     func(ctx context.Context, resources []resource.ID)
	OnAllocErrorF   func(ctx context.Context, n int, err error)
	BeforeInitF     func(ctx context.Context, resources []resource.Resource) error
	AfterInitF      func(ctx context.Context, resources []resource.Resource)
	OnInitErrorF    func(ctx context.Context, resources []resource.Resource, err error)
	BeforeDeallocF  func(ctx context.Context, resources []resource.ID) error
	AfterDeallocF   func(ctx context.Context, resources []resource.ID)
	OnDeallocErrorF func(ctx context.Context, resources []resource.ID, err error)
}

func (observer FunctionalObserver) BeforeLoad(ctx context.Context, ids []resource.ID) error {
	if observer.BeforeLoadF == nil {
		return nil
	}
	return observer.BeforeLoadF(ctx, ids)
}

func (observer FunctionalObserver) AfterLoad(ctx context.Context, resources []resource.Resource) {
	if observer.AfterLoadF == nil {
		return
	}
	observer.AfterLoadF(ctx, resources)
}

func (observer FunctionalObserver) OnLoadError(ctx context.Context, ids []resource.ID, err error) {
	if observer.OnLoadErrorF == nil {
		return
	}
	observer.OnLoadErrorF(ctx, ids, err)
}

func (observer FunctionalObserver) BeforeAlloc(ctx context.Context, n int) error {
	if observer.BeforeAllocF == nil {
		return nil
	}
	return observer.BeforeAllocF(ctx, n)
}

func (observer FunctionalObserver) AfterAlloc(ctx context.Context, resources []resource.ID) {
	if observer.AfterAllocF == nil {
		return
	}
	observer.AfterAllocF(ctx, resources)
}

func (observer FunctionalObserver) OnAllocError(ctx context.Context, n int, err error) {
	if observer.OnAllocErrorF == nil {
		return
	}
	observer.OnAllocErrorF(ctx, n, err)
}

func (observer FunctionalObserver) BeforeInit(ctx context.Context, resources []resource.Resource) error {
	if observer.BeforeInitF == nil {
		return nil
	}
	return observer.BeforeInitF(ctx, resources)
}

func (observer FunctionalObserver) AfterInit(ctx context.Context, resources []resource.Resource) {
	if observer.AfterInitF == nil {
		return
	}
	observer.AfterInitF(ctx, resources)
}

func (observer FunctionalObserver) OnInitError(ctx context.Context, resources []resource.Resource, err error) {
	if observer.OnInitErrorF == nil {
		return
	}
	observer.OnInitErrorF(ctx, resources, err)
}

func (observer FunctionalObserver) BeforeDealloc(ctx context.Context, resources []resource.ID) error {
	if observer.BeforeDeallocF == nil {
		return nil
	}
	return observer.BeforeDeallocF(ctx, resources)
}

func (observer FunctionalObserver) AfterDealloc(ctx context.Context, resources []resource.ID) {
	if observer.AfterDeallocF == nil {
		return
	}
	observer.AfterDeallocF(ctx, resources)
}

func (observer FunctionalObserver) OnDeallocError(ctx context.Context, resources []resource.ID, err error) {
	if observer.OnDeallocErrorF == nil {
		return
	}
	observer.OnDeallocErrorF(ctx, resources, err)
}

type DummyObserver struct{}

func (observer DummyObserver) BeforeLoad(ctx context.Context, ids []resource.ID) error {
	return nil
}

func (observer DummyObserver) AfterLoad(ctx context.Context, resources []resource.Resource) {}

func (observer DummyObserver) OnLoadError(ctx context.Context, ids []resource.ID, err error) {}

func (observer DummyObserver) BeforeAlloc(ctx context.Context, n int) error {
	return nil
}

func (observer DummyObserver) AfterAlloc(ctx context.Context, resources []resource.ID) {}

func (observer DummyObserver) OnAllocError(ctx context.Context, n int, err error) {}

func (observer DummyObserver) BeforeInit(ctx context.Context, resources []resource.Resource) error {
	return nil
}

func (observer DummyObserver) AfterInit(ctx context.Context, resources []resource.Resource) {}

func (observer DummyObserver) OnInitError(ctx context.Context, resources []resource.Resource, err error) {
}

func (observer DummyObserver) BeforeDealloc(ctx context.Context, resources []resource.ID) error {
	return nil
}

func (observer DummyObserver) AfterDealloc(ctx context.Context, resources []resource.ID) {}

func (observer DummyObserver) OnDeallocError(ctx context.Context, resources []resource.ID, err error) {
}

type MultiObserver []Observer

func (observers MultiObserver) BeforeLoad(ctx context.Context, ids []resource.ID) error {
	for _, observer := range observers {
		if err := observer.BeforeLoad(ctx, ids); err != nil {
			return err
		}
	}
	return nil
}

func (observers MultiObserver) AfterLoad(ctx context.Context, resources []resource.Resource) {
	for _, observer := range observers {
		observer.AfterLoad(ctx, resources)
	}
}

func (observers MultiObserver) OnLoadError(ctx context.Context, ids []resource.ID, err error) {
	for _, observer := range observers {
		observer.OnLoadError(ctx, ids, err)
	}
}

func (observers MultiObserver) BeforeAlloc(ctx context.Context, n int) error {
	for _, observer := range observers {
		if err := observer.BeforeAlloc(ctx, n); err != nil {
			return err
		}
	}
	return nil
}

func (observers MultiObserver) AfterAlloc(ctx context.Context, resources []resource.ID) {
	for _, observer := range observers {
		observer.AfterAlloc(ctx, resources)
	}
}

func (observers MultiObserver) OnAllocError(ctx context.Context, n int, err error) {
	for _, observer := range observers {
		observer.OnAllocError(ctx, n, err)
	}
}

func (observers MultiObserver) BeforeInit(ctx context.Context, resources []resource.Resource) error {
	for _, observer := range observers {
		if err := observer.BeforeInit(ctx, resources); err != nil {
			return err
		}
	}
	return nil
}

func (observers MultiObserver) AfterInit(ctx context.Context, resources []resource.Resource) {
	for _, observer := range observers {
		observer.AfterInit(ctx, resources)
	}
}

func (observers MultiObserver) OnInitError(ctx context.Context, resources []resource.Resource, err error) {
	for _, observer := range observers {
		observer.OnInitError(ctx, resources, err)
	}
}

func (observers MultiObserver) BeforeDealloc(ctx context.Context, resources []resource.ID) error {
	for _, observer := range observers {
		if err := observer.BeforeDealloc(ctx, resources); err != nil {
			return err
		}
	}
	return nil
}

func (observers MultiObserver) AfterDealloc(ctx context.Context, resources []resource.ID) {
	for _, observer := range observers {
		observer.AfterDealloc(ctx, resources)
	}
}

func (observers MultiObserver) OnDeallocError(ctx context.Context, resources []resource.ID, err error) {
	for _, observer := range observers {
		observer.OnDeallocError(ctx, resources, err)
	}
}

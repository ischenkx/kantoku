package resources

import (
	"context"
	"fmt"
	"github.com/ischenkx/kantoku/pkg/core"
)

type Observer interface {
	BeforeLoad(ctx context.Context, ids []string) error
	AfterLoad(ctx context.Context, resources []core.Resource)
	OnLoadError(ctx context.Context, ids []string, err error)

	BeforeAlloc(ctx context.Context, n int) error
	AfterAlloc(ctx context.Context, resources []string)
	OnAllocError(ctx context.Context, n int, err error)

	BeforeInit(ctx context.Context, resources []core.Resource) error
	AfterInit(ctx context.Context, resources []core.Resource)
	OnInitError(ctx context.Context, resources []core.Resource, err error)

	BeforeDealloc(ctx context.Context, resources []string) error
	AfterDealloc(ctx context.Context, resources []string)
	OnDeallocError(ctx context.Context, resources []string, err error)
}

func Observe(raw core.ResourceDB, observer Observer) *Observable {
	return &Observable{
		raw:      raw,
		observer: observer,
	}
}

type Observable struct {
	raw      core.ResourceDB
	observer Observer
}

func (storage *Observable) Load(ctx context.Context, ids ...string) ([]core.Resource, error) {
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

func (storage *Observable) Alloc(ctx context.Context, amount int) ([]string, error) {
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

func (storage *Observable) Init(ctx context.Context, resources []core.Resource) error {
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

func (storage *Observable) Dealloc(ctx context.Context, ids []string) error {
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
	BeforeLoadF     func(ctx context.Context, ids []string) error
	AfterLoadF      func(ctx context.Context, resources []core.Resource)
	OnLoadErrorF    func(ctx context.Context, ids []string, err error)
	BeforeAllocF    func(ctx context.Context, n int) error
	AfterAllocF     func(ctx context.Context, resources []string)
	OnAllocErrorF   func(ctx context.Context, n int, err error)
	BeforeInitF     func(ctx context.Context, resources []core.Resource) error
	AfterInitF      func(ctx context.Context, resources []core.Resource)
	OnInitErrorF    func(ctx context.Context, resources []core.Resource, err error)
	BeforeDeallocF  func(ctx context.Context, resources []string) error
	AfterDeallocF   func(ctx context.Context, resources []string)
	OnDeallocErrorF func(ctx context.Context, resources []string, err error)
}

func (observer FunctionalObserver) BeforeLoad(ctx context.Context, ids []string) error {
	if observer.BeforeLoadF == nil {
		return nil
	}
	return observer.BeforeLoadF(ctx, ids)
}

func (observer FunctionalObserver) AfterLoad(ctx context.Context, resources []core.Resource) {
	if observer.AfterLoadF == nil {
		return
	}
	observer.AfterLoadF(ctx, resources)
}

func (observer FunctionalObserver) OnLoadError(ctx context.Context, ids []string, err error) {
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

func (observer FunctionalObserver) AfterAlloc(ctx context.Context, resources []string) {
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

func (observer FunctionalObserver) BeforeInit(ctx context.Context, resources []core.Resource) error {
	if observer.BeforeInitF == nil {
		return nil
	}
	return observer.BeforeInitF(ctx, resources)
}

func (observer FunctionalObserver) AfterInit(ctx context.Context, resources []core.Resource) {
	if observer.AfterInitF == nil {
		return
	}
	observer.AfterInitF(ctx, resources)
}

func (observer FunctionalObserver) OnInitError(ctx context.Context, resources []core.Resource, err error) {
	if observer.OnInitErrorF == nil {
		return
	}
	observer.OnInitErrorF(ctx, resources, err)
}

func (observer FunctionalObserver) BeforeDealloc(ctx context.Context, resources []string) error {
	if observer.BeforeDeallocF == nil {
		return nil
	}
	return observer.BeforeDeallocF(ctx, resources)
}

func (observer FunctionalObserver) AfterDealloc(ctx context.Context, resources []string) {
	if observer.AfterDeallocF == nil {
		return
	}
	observer.AfterDeallocF(ctx, resources)
}

func (observer FunctionalObserver) OnDeallocError(ctx context.Context, resources []string, err error) {
	if observer.OnDeallocErrorF == nil {
		return
	}
	observer.OnDeallocErrorF(ctx, resources, err)
}

type DummyObserver struct{}

func (observer DummyObserver) BeforeLoad(ctx context.Context, ids []string) error {
	return nil
}

func (observer DummyObserver) AfterLoad(ctx context.Context, resources []core.Resource) {}

func (observer DummyObserver) OnLoadError(ctx context.Context, ids []string, err error) {}

func (observer DummyObserver) BeforeAlloc(ctx context.Context, n int) error {
	return nil
}

func (observer DummyObserver) AfterAlloc(ctx context.Context, resources []string) {}

func (observer DummyObserver) OnAllocError(ctx context.Context, n int, err error) {}

func (observer DummyObserver) BeforeInit(ctx context.Context, resources []core.Resource) error {
	return nil
}

func (observer DummyObserver) AfterInit(ctx context.Context, resources []core.Resource) {}

func (observer DummyObserver) OnInitError(ctx context.Context, resources []core.Resource, err error) {
}

func (observer DummyObserver) BeforeDealloc(ctx context.Context, resources []string) error {
	return nil
}

func (observer DummyObserver) AfterDealloc(ctx context.Context, resources []string) {}

func (observer DummyObserver) OnDeallocError(ctx context.Context, resources []string, err error) {
}

package mempool

import (
	"context"
	"math/rand"
	"sync"
	"time"
)

func New[T any]() *Pool[T] {
	return &Pool[T]{}
}

type Pool[T any] struct {
	buffer  []T
	readers []chan T
	perm    []int
	closed  bool
	mu      sync.RWMutex
}

func (p *Pool[T]) Read(ctx context.Context) (<-chan T, error) {
	p.mu.Lock()
	defer p.mu.Unlock()

	channel := make(chan T, 128)

	go func() {
		<-ctx.Done()
		p.mu.Lock()
		defer p.mu.Unlock()
		for i, reader := range p.readers {
			if reader == channel {
				p.readers[i], p.readers[len(p.readers)-1] = p.readers[len(p.readers)-1], p.readers[i]
				p.readers = p.readers[:len(p.readers)-1]
				break
			}
		}
	}()

	p.readers = append(p.readers, channel)

	return channel, nil
}

func (p *Pool[T]) Write(ctx context.Context, item T) error {
	p.mu.Lock()
	defer p.mu.Unlock()

	if !p.write(item) {
		p.buffer = append(p.buffer, item)
	}

	return nil
}

func (p *Pool[T]) Close() {
	p.mu.Lock()
	defer p.mu.Unlock()
	for _, channel := range p.readers {
		close(channel)
	}

	p.closed = true
	p.buffer = nil
	p.readers = nil
}

func (p *Pool[T]) runFlusher() {
	ticker := time.NewTicker(time.Second * 3)
	defer ticker.Stop()

	for range ticker.C {
		if !p.flush() {
			break
		}
	}
}

func (p *Pool[T]) flush() bool {
	p.mu.Lock()
	defer p.mu.Unlock()

	if p.closed {
		return false
	}

	for len(p.buffer) != 0 {
		index := len(p.buffer) - 1
		if !p.write(p.buffer[index]) {
			break
		}
		p.buffer = p.buffer[:index]
	}

	return true
}

func (p *Pool[T]) write(item T) bool {
	for _, index := range p.permute() {
		reader := p.readers[index]
		select {
		case reader <- item:
			return true
		default:
		}
	}
	return false
}

func (p *Pool[T]) permute() []int {
	if len(p.perm) > len(p.readers) {
		for i := 0; i < len(p.readers); i++ {
			p.perm[i] = i
		}
		p.perm = p.perm[:len(p.readers)]
	} else {
		for i := len(p.perm); i < len(p.readers); i++ {
			p.perm = append(p.perm, i)
		}
	}

	rand.Shuffle(len(p.perm), func(i, j int) { p.perm[i], p.perm[j] = p.perm[j], p.perm[i] })

	return p.perm
}

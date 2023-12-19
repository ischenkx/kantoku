package mempool

import (
	"context"
	"go/types"
	"kantoku/common/data/pool"
	"kantoku/common/data/transactional"
	"math/rand"
	"sync"
	"time"
)

var _ pool.Pool[types.Object] = &Pool[types.Object]{}

func New[T any](config Config) *Pool[T] {
	p := &Pool[T]{config: config}
	go p.runFlusher()
	return p
}

type Pool[T any] struct {
	buffer  []T
	readers []chan transactional.Object[T]
	perm    []int
	closed  bool
	mu      sync.RWMutex
	config  Config
}

type Config struct {
	BufferSize  int
	FlushPeriod time.Duration
}

var DefaultConfig = Config{
	BufferSize:  0,
	FlushPeriod: time.Second,
}

func (p *Pool[T]) Read(ctx context.Context) (<-chan transactional.Object[T], error) {
	p.mu.Lock()
	defer p.mu.Unlock()

	channel := make(chan transactional.Object[T], p.config.BufferSize)

	go func() {
		<-ctx.Done()
		p.mu.Lock()
		defer p.mu.Unlock()
		for i, reader := range p.readers {
			if reader == channel {
				close(channel)
				p.readers[i], p.readers[len(p.readers)-1] = p.readers[len(p.readers)-1], p.readers[i]
				p.readers = p.readers[:len(p.readers)-1]
				break
			}
		}
	}()

	p.readers = append(p.readers, channel)

	return channel, nil
}

func (p *Pool[T]) Write(_ context.Context, items ...T) error {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.buffer = append(p.buffer, items...)

	p.lockedFlush() // doing that to guarantee no delay if someone is ready to read
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
	ticker := time.NewTicker(p.config.FlushPeriod)
	defer ticker.Stop()

	for range ticker.C {
		if !p.flush() {
			break
		}
	}
}

func (p *Pool[T]) lockedFlush() bool {
	if p.closed {
		return false
	}

	for len(p.buffer) != 0 {
		if !p.write() {
			break
		}
	}

	return true
}

func (p *Pool[T]) flush() bool {
	p.mu.Lock()
	defer p.mu.Unlock()

	return p.lockedFlush()
}

func (p *Pool[T]) write() bool {
	for _, index := range p.permute() {
		reader := p.readers[index]
		statusChan := make(chan bool)
		select {
		case reader <- &Transaction[T]{data: p.buffer[0], success: statusChan}:
			// might want to add context with deadline
			success := <-statusChan
			if success {
				p.buffer = p.buffer[1:]
			}

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

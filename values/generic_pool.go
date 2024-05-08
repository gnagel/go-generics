package values

import (
	"sync"
	"sync/atomic"
)

type ValuePool[T any] struct {
	pool     sync.Pool
	len, cap atomic.Int64
	reset    Reset[T]
}

func NewValuePool[T any](alloc Alloc[T], reset Reset[T]) *ValuePool[T] {
	if nil == alloc {
		alloc = DefaultAlloc[T](nil)
	}
	if nil == reset {
		reset = DefaultReset[T]
	}

	output := &ValuePool[T]{reset: reset}
	output.pool.New = valueAlloc(&output.cap, alloc)
	return output
}

func (p *ValuePool[T]) Get() T {
	go func(counter *atomic.Int64) { counter.Add(1) }(&p.len)
	return p.pool.Get().(T)
}

func (p *ValuePool[T]) Put(value T) {
	value = p.reset(value)
	p.pool.Put(value)
	go func(counter *atomic.Int64) { counter.Add(-1) }(&p.len)
}

func (p *ValuePool[T]) Len() int64 {
	return p.len.Load()
}

func (p *ValuePool[T]) Cap() int64 {
	return p.cap.Load()
}

func valueAlloc[T any](cap *atomic.Int64, alloc Alloc[T]) func() any {
	return func() any {
		go func() { cap.Add(1) }()
		return alloc()
	}
}

package slices

import (
	"sync"
	"sync/atomic"
)

type SlicePool[T any] struct {
	pool     sync.Pool
	len, cap atomic.Int64
	reset    ResetSlice[T]
}

func NewSlicePool[T any](alloc AllocSlice[T], reset ResetSlice[T], defaultLen, defaultCap int) *SlicePool[T] {
	output := &SlicePool[T]{reset: reset}
	output.pool.New = sliceAlloc(&output.cap, alloc, defaultLen, defaultCap)
	return output
}

func (p *SlicePool[T]) Get() []T {
	go func(counter *atomic.Int64) { counter.Add(1) }(&p.len)
	return p.pool.Get().([]T)
}

func (p *SlicePool[T]) Put(value []T) {
	value = p.reset(value)
	p.pool.Put(value)
	go func(counter *atomic.Int64) { counter.Add(-1) }(&p.len)
}

func (p *SlicePool[T]) Len() int64 {
	return p.len.Load()
}

func (p *SlicePool[T]) Cap() int64 {
	return p.cap.Load()
}

func sliceAlloc[T any](cap *atomic.Int64, alloc AllocSlice[T], defaultLen, defaultCap int) func() any {
	return func() any {
		go func() { cap.Add(1) }()
		return alloc(defaultLen, defaultCap)
	}
}

package pools

import (
	"github.com/go-generics-playground/generics/functions"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
	"log"
	"sync"
	"sync/atomic"
)

type SlicePool[T any] struct {
	pool     sync.Pool
	len, cap atomic.Int64
	reset    functions.ResetSlice[T]
}

func NewSlicePool[T any](alloc functions.AllocSlice[T], reset functions.ResetSlice[T], defaultLen, defaultCap int) *SlicePool[T] {
	if nil == alloc {
		log.Panic("alloc is required for ValuePool")
	}
	if nil == reset {
		log.Panic("reset is required for ValuePool")
	}
	if defaultLen > defaultCap {
		log.Panicf("len(%d) must be <= cap(%d)", defaultLen, defaultCap)
	}
	output := &SlicePool[T]{reset: reset}
	output.pool.New = func() any {
		go func() { output.cap.Add(1) }()
		return alloc(defaultLen, defaultCap)
	}
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

type ProtoSlicePool struct {
	*SlicePool[proto.Message]
}

func NewProtoSlicePool(descriptor protoreflect.Message, defaultLen, defaultCap int) *ProtoSlicePool {
	alloc := functions.AllocProtoSlice(descriptor)
	reset := functions.ResetProtoSlice(descriptor)

	return &ProtoSlicePool{
		SlicePool: NewSlicePool[proto.Message](alloc, reset, defaultLen, defaultCap),
	}
}

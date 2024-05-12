package pools

import (
	"github.com/go-generics-playground/generics/functions"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
	"log"
	"sync"
	"sync/atomic"
)

type ValuePool[T any] struct {
	pool     sync.Pool
	len, cap atomic.Int64
	reset    functions.Reset[T]
}

func NewValuePool[T any](alloc functions.Alloc[T], reset functions.Reset[T]) *ValuePool[T] {
	if nil == alloc {
		log.Panic("alloc is required for ValuePool")
	}
	if nil == reset {
		log.Panic("reset is required for ValuePool")
	}

	output := &ValuePool[T]{reset: reset}
	output.pool.New = func() any {
		go func() { output.cap.Add(1) }()
		return alloc()
	}
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

type ProtoPool struct {
	*ValuePool[proto.Message]
}

func NewProtoPool(descriptor protoreflect.Message) *ProtoPool {
	alloc := functions.AllocProto(descriptor)
	reset := functions.ResetProto(descriptor)

	return &ProtoPool{
		ValuePool: NewValuePool[proto.Message](alloc, reset),
	}
}

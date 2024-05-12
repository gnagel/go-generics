package pools

import (
	"google.golang.org/protobuf/proto"
)

// SyncPool is generic wrapper around a sync.Pool with some additional helper methods for working with generics.
// All the pools in the nested packages implement this base interface:
// *values
// *slices
type SyncPool[T any] interface {
	Get() T
	Put(value T)
	Len() int64
	Cap() int64
}

var _ SyncPool[any] = &ValuePool[any]{}
var _ SyncPool[[]any] = &SlicePool[any]{}

var _ SyncPool[proto.Message] = &ProtoPool{}
var _ SyncPool[[]proto.Message] = &ProtoSlicePool{}

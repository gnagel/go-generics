package functions

import (
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
)

// Dealloc releases a value back to the garbage collector and releases any internals.
type Dealloc[T any] func(T)

// DeallocSlice is a generic function for releasing an existing slice of type T.
// This is the same as slices.Clip, except with the added call to dealloc for each element
// This returns an empty slice of length 0 & capacity 0
type DeallocSlice[T any] func([]T) []T

var _ Dealloc[any] = DefaultDealloc[any]
var _ DeallocSlice[any] = DefaultDeallocSlice[any](nil)
var _ Dealloc[proto.Message] = DeallocProto(exampleDescriptor)
var _ DeallocSlice[proto.Message] = DeallocProtoSlice(exampleDescriptor)

// DefaultDealloc implements the Dealloc[T] interface for generic values
func DefaultDealloc[T any](_ T) {}

// DefaultDeallocSlice implements the DeallocSlice[T] interface for generic values
func DefaultDeallocSlice[T any](dealloc Dealloc[T]) DeallocSlice[T] {
	if nil == dealloc {
		dealloc = DefaultDealloc[T]
	}

	return func(input []T) []T {
		for index, value := range input {
			dealloc(value)
			var t T
			input[index] = t
		}
		return input[:0:0]
	}
}

// DeallocProto implements the Dealloc[T] interface for a protobuf descriptor
func DeallocProto(_ protoreflect.Message) Dealloc[proto.Message] {
	return func(msg proto.Message) {
		proto.Reset(msg)
	}
}

// DeallocProtoSlice implements the DeallocSlice[T] interface for a protobuf descriptor
func DeallocProtoSlice(descriptor protoreflect.Message) DeallocSlice[proto.Message] {
	dealloc := DeallocProto(descriptor)

	return func(input []proto.Message) []proto.Message {
		for index, value := range input {
			dealloc(value)
			input[index] = nil
		}
		return input[:0:0]
	}
}

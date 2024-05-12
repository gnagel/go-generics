package functions

import (
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
)

// Reset zeros out a value and resets it back to zero state.
type Reset[T any] func(T) T

// ResetSlice is a generic function for resetting an existing slice of type T.
// This returns an empty slice of length 0, but the capacity can be > 0
type ResetSlice[T any] func([]T) []T

var _ Reset[any] = DefaultReset[any]
var _ ResetSlice[any] = DefaultResetSlice[any](nil)
var _ Reset[proto.Message] = ResetProto(exampleDescriptor)
var _ ResetSlice[proto.Message] = ResetProtoSlice(exampleDescriptor)

// DefaultReset implements the Reset[T] interface for generic values
func DefaultReset[T any](t T) T {
	var v T
	t = v
	return t
}

// DefaultResetSlice implements the ResetSlice[T] interface for generic values
func DefaultResetSlice[T any](reset Reset[T]) ResetSlice[T] {
	if nil == reset {
		reset = DefaultReset[T]
	}
	return func(input []T) []T {
		for index, value := range input {
			input[index] = reset(value)
		}
		return input[:0:cap(input)]
	}
}

// ResetProto implements the Reset[T] interface for a protobuf descriptor
func ResetProto(descriptor protoreflect.Message) Reset[proto.Message] {
	return func(msg proto.Message) proto.Message {
		proto.Reset(msg)
		return msg
	}
}

// ResetProtoSlice implements the ResetSlice[T] interface for a protobuf descriptor
func ResetProtoSlice(descriptor protoreflect.Message) ResetSlice[proto.Message] {
	reset := ResetProto(descriptor)
	return func(input []proto.Message) []proto.Message {
		for index, value := range input {
			input[index] = reset(value)
		}
		return input[:0]
	}
}

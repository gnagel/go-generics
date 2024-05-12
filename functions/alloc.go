package functions

import (
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/types/known/structpb"
)

// Alloc is a generic function for creating a new instance of type T.
type Alloc[T any] func() T

// AllocSlice is a generic function for creating a new slice of type T.
type AllocSlice[T any] func(defaultLen, defaultCap int) []T

var _ Alloc[any] = DefaultAlloc[any](nil)
var _ AllocSlice[any] = DefaultAllocSlice[any](nil)
var _ Alloc[proto.Message] = AllocProto(exampleDescriptor)
var _ AllocSlice[proto.Message] = AllocProtoSlice(exampleDescriptor)

// DefaultAlloc implements the Alloc[T] interface for generic values
func DefaultAlloc[T any](defaultValue T) Alloc[T] {
	return func() T {
		var t = defaultValue
		return t
	}
}

// DefaultAllocSlice implements the AllocSlice[T] interface for generic values
func DefaultAllocSlice[T any](alloc Alloc[T]) AllocSlice[T] {
	if nil == alloc {
		var t T
		alloc = DefaultAlloc[T](t)
	}

	return func(defaultLen, defaultCap int) []T {
		output := make([]T, defaultLen, defaultCap)
		for index := range output[:defaultLen] {
			output[index] = alloc()
		}
		return output
	}
}

// AllocProto implements the Alloc[T] interface for a protobuf descriptor
func AllocProto(descriptor protoreflect.Message) Alloc[proto.Message] {
	return func() proto.Message {
		var msg protoreflect.ProtoMessage
		msg = descriptor.New().(protoreflect.ProtoMessage)
		return msg
	}
}

// AllocProtoSlice implements the AllocSlice[T] interface for a protobuf descriptor
func AllocProtoSlice(descriptor protoreflect.Message) AllocSlice[proto.Message] {
	alloc := AllocProto(descriptor)
	return func(defaultLen, defaultCap int) []proto.Message {
		output := make([]proto.Message, defaultLen, defaultCap)
		for index := range output[:defaultLen] {
			output[index] = alloc()
		}
		return output
	}
}

// Sample proto & descriptor
var exampleMsg structpb.Struct
var exampleDescriptor = exampleMsg.ProtoReflect()

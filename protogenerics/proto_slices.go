package protogenerics

import (
	"github.com/go-generics-playground/generics/slices"
	"google.golang.org/protobuf/proto"
)
import "google.golang.org/protobuf/reflect/protoreflect"

// AllocSlice, DeallocSlice, and ResetSlice all implement the generic interface
var _ slices.AllocSlice[proto.Message] = ProtoSliceAlloc(exampleDescriptor)
var _ slices.DeallocSlice[proto.Message] = ProtoSliceDealloc(exampleDescriptor)
var _ slices.ResetSlice[proto.Message] = ProtoSliceReset(exampleDescriptor)

// ProtoSliceAlloc implements the AllocSlice[T] interface for a protobuf descriptor
func ProtoSliceAlloc(descriptor protoreflect.Message) slices.AllocSlice[proto.Message] {
	alloc := ProtoAlloc(descriptor)
	return func(defaultLen, defaultCap int) []proto.Message {
		output := make([]proto.Message, defaultLen, defaultCap)
		for index := range output[:defaultLen] {
			output[index] = alloc()
		}
		return output
	}
}

// ProtoSliceDealloc implements the DeallocSlice[T] interface for a protobuf descriptor
func ProtoSliceDealloc(descriptor protoreflect.Message) slices.DeallocSlice[proto.Message] {
	dealloc := ProtoDealloc(descriptor)

	return func(input []proto.Message) []proto.Message {
		for index, value := range input {
			dealloc(value)
			input[index] = nil
		}
		return input[:0:0]
	}
}

// ProtoSliceReset implements the ResetSlice[T] interface for a protobuf descriptor
func ProtoSliceReset(descriptor protoreflect.Message) slices.ResetSlice[proto.Message] {
	reset := ProtoReset(descriptor)
	return func(input []proto.Message) []proto.Message {
		for index, value := range input {
			input[index] = reset(value)
		}
		return input[:0]
	}
}

package values

import (
	"github.com/go-generics-playground/generics/values"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/structpb"
	"slices"
)
import "google.golang.org/protobuf/reflect/protoreflect"

var exampleMsg structpb.Struct
var exampleDescriptor = exampleMsg.ProtoReflect()

// AllocSlice, DeallocSlice, and ResetSlice all implement the generic interface
var _ AllocSlice[proto.Message] = ProtoSliceAlloc(exampleDescriptor)
var _ DeallocSlice[proto.Message] = ProtoSliceDealloc(exampleDescriptor)
var _ ResetSlice[proto.Message] = ProtoSliceReset(exampleDescriptor)

// ProtoAlloc implements the Alloc[T] interface for a protobuf descriptor
var ProtoAlloc = values.ProtoAlloc

// ProtoDealloc implements the Dealloc[T] interface for a protobuf descriptor
var ProtoDealloc = values.ProtoDealloc

// ProtoReset implements the Reset[T] interface for a protobuf descriptor
var ProtoReset = values.ProtoReset

// ProtoSliceAlloc implements the AllocSlice[T] interface for a protobuf descriptor
func ProtoSliceAlloc(descriptor protoreflect.Message) AllocSlice[proto.Message] {
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
func ProtoSliceDealloc(descriptor protoreflect.Message) DeallocSlice[proto.Message] {
	dealloc := ProtoDealloc(descriptor)

	return func(input []proto.Message) []proto.Message {
		for index, value := range input {
			dealloc(value)
			input[index] = nil
		}
		return slices.Clip(input)
	}
}

// ProtoSliceReset implements the ResetSlice[T] interface for a protobuf descriptor
func ProtoSliceReset(descriptor protoreflect.Message) ResetSlice[proto.Message] {
	reset := ProtoReset(descriptor)
	return func(input []proto.Message) []proto.Message {
		for index, value := range input {
			input[index] = reset(value)
		}
		return input[:0]
	}
}

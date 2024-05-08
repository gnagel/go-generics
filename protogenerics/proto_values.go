package protogenerics

import (
	"github.com/go-generics-playground/generics/values"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/structpb"
)
import "google.golang.org/protobuf/reflect/protoreflect"

var exampleMsg structpb.Struct
var exampleDescriptor = exampleMsg.ProtoReflect()

// Alloc, Dealloc, and Reset all implement the generic interface
var _ values.Alloc[proto.Message] = ProtoAlloc(exampleDescriptor)
var _ values.Dealloc[proto.Message] = ProtoDealloc(exampleDescriptor)
var _ values.Reset[proto.Message] = ProtoReset(exampleDescriptor)

// ProtoAlloc implements the Alloc[T] interface for a protobuf descriptor
func ProtoAlloc(descriptor protoreflect.Message) values.Alloc[proto.Message] {
	return func() proto.Message {
		var msg protoreflect.ProtoMessage
		msg = descriptor.New().(protoreflect.ProtoMessage)
		return msg
	}
}

// ProtoDealloc implements the Dealloc[T] interface for a protobuf descriptor
func ProtoDealloc(descriptor protoreflect.Message) values.Dealloc[proto.Message] {
	return func(msg proto.Message) {
		proto.Reset(msg)
	}
}

// ProtoReset implements the Reset[T] interface for a protobuf descriptor
func ProtoReset(descriptor protoreflect.Message) values.Reset[proto.Message] {
	return func(msg proto.Message) proto.Message {
		proto.Reset(msg)
		return msg
	}
}

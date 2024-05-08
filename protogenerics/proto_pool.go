package protogenerics

import (
	"github.com/go-generics-playground/generics/slices"
	"github.com/go-generics-playground/generics/values"
	"google.golang.org/protobuf/proto"
)
import "google.golang.org/protobuf/reflect/protoreflect"

type ProtoPool struct {
	*values.ValuePool[proto.Message]
}

func NewProtoPool(descriptor protoreflect.Message) *ProtoPool {
	alloc := ProtoAlloc(descriptor)
	reset := ProtoReset(descriptor)

	return &ProtoPool{
		ValuePool: values.NewValuePool[proto.Message](alloc, reset),
	}
}

type ProtoSlicePool struct {
	*slices.SlicePool[proto.Message]
}

func NewProtoSlicePool(descriptor protoreflect.Message, defaultLen, defaultCap int) *ProtoSlicePool {
	alloc := ProtoSliceAlloc(descriptor)
	reset := ProtoSliceReset(descriptor)

	return &ProtoSlicePool{
		SlicePool: slices.NewSlicePool[proto.Message](alloc, reset, defaultLen, defaultCap),
	}
}

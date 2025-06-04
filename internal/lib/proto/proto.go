package protolib

import (
	"google.golang.org/protobuf/types/known/structpb"
)

func MapToProtobufStruct(m map[string]interface{}) (*structpb.Struct, error) {
	return structpb.NewStruct(m)
}

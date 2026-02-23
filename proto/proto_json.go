package proto

import (
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
)

func MarshalProto(proto proto.Message) (string, error) {
	opts := protojson.MarshalOptions{
		Indent: "  ",
	}
	json, err := opts.Marshal(proto)
	if err != nil {
		return "", err
	}
	return string(json), nil
}

package message

import (
	"encoding/json"
	"errors"
	"google.golang.org/protobuf/proto"
)

var ProtoBuffCodec = protobufCodec{}
var JsonCodec = jsonCodec{}
var DefaultCodec = ProtoBuffCodec

func UnmarshallJson(json_ string, i interface{}) error {
	return json.Unmarshal([]byte(json_), i)
}

type Codec interface {
	Decode(data []byte, i interface{}) error
	Encode(i interface{}) ([]byte, error)
}

type protobufCodec struct {
}

func (p protobufCodec) Decode(data []byte, i interface{}) error {
	message, ok := i.(proto.Message)
	if !ok {
		return errors.New("illegal argument, not implement proto.Message")
	}
	return proto.Unmarshal(data, message)
}

func (p protobufCodec) Encode(i interface{}) ([]byte, error) {
	message, ok := i.(proto.Message)
	if !ok {
		return nil, errors.New("illegal argument, not implement proto.Message")
	}
	return proto.Marshal(message)
}

type jsonCodec struct {
}

func (j jsonCodec) Decode(data []byte, i interface{}) error {
	return json.Unmarshal(data, i)
}

func (j jsonCodec) Encode(i interface{}) ([]byte, error) {
	return json.Marshal(i)
}

package rpc

import (
	"fmt"
	ggProto "github.com/gogo/protobuf/proto"
	"github.com/smallnest/rpcx/protocol"
	"github.com/smallnest/rpcx/share"
	"go_im/service/pb"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/anypb"
)

//var SerialTypeProtoBuffWrapAny = protocol.SerializeType(100)

const (
	_ = protocol.Thrift + 1
	SerialTypeProtoBuffWrapAny
)

type PB4AnyWrapperCodec struct{}

func init() {
	share.Codecs[SerialTypeProtoBuffWrapAny] = &PB4AnyWrapperCodec{}
}

func (c PB4AnyWrapperCodec) Encode(i interface{}) ([]byte, error) {

	if m, ok := i.(ggProto.Marshaler); ok {
		return m.Marshal()
	}

	if m, ok := i.(proto.Message); ok {

		a, ok2 := i.(*pb.RouteReqParam)
		if ok2 {
			return proto.Marshal(a.GetData())
		}
		return proto.Marshal(m)
	}

	return nil, fmt.Errorf("%T is not a pb.Message", i)
}

func (c PB4AnyWrapperCodec) Decode(data []byte, i interface{}) error {

	if m, ok := i.(ggProto.Unmarshaler); ok {
		return m.Unmarshal(data)
	}

	if m, ok := i.(proto.Message); ok {
		any := &anypb.Any{}
		e := proto.Unmarshal(data, any)
		if e != nil && any.MessageIs(m) {
			return any.UnmarshalTo(m)
		}
		return proto.Unmarshal(data, m)
	}

	return fmt.Errorf("%T is not a pb.Message", i)
}

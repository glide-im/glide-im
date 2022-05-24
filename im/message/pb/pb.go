package pb

import (
	"github.com/glide-im/glideim/pkg/logger"
	"github.com/glide-im/glideim/protobuf/gen/pb_im"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/anypb"
	"reflect"
)

func NewMessage(seq int64, action string, data interface{}) *pb_im.CommMessage {
	message := &pb_im.CommMessage{
		Ver:    0,
		Seq:    seq,
		Action: action,
		Data:   nil,
	}
	if data == nil {
		return message
	}
	p, ok := data.(proto.Message)
	if !ok {
		logger.E("%s is not proto.Message", reflect.TypeOf(data).String())
		return message
	}
	any, err := anypb.New(p)
	if err != nil {
		logger.E("marshal pb message data error %v", err)
	}
	message.Data = any
	return message
}

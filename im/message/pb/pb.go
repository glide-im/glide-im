package pb

import (
	"go_im/pkg/logger"
	"go_im/protobuff/pb_im"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/anypb"
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
		logger.E("%v is not proto.Message")
		return message
	}
	any, err := anypb.New(p)
	if err != nil {
		logger.E("marshal pb message data error %v", err)
	}
	message.Data = any
	return message
}

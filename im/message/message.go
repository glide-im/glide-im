package message

import (
	"go_im/protobuf/gen/pb_im"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/anypb"
)

// ChatMessage 上行消息, 表示服务端收到发送者的消息
type ChatMessage struct {
	*pb_im.ChatMessage
}

func NewChatMessage(mid, seq, from, to int64, typ int32, content string, sendAt int64) ChatMessage {
	return ChatMessage{
		&pb_im.ChatMessage{
			Mid:     mid,
			Seq:     seq,
			From:    from,
			To:      to,
			Type:    typ,
			Content: content,
			SendAt:  sendAt,
		},
	}
}

// AckMessage 服务端通知发送者的服务端收到消息
type AckMessage struct {
	*pb_im.AckMessage
}

func NewAckMessage(mid int64, seq int64) AckMessage {
	return AckMessage{
		&pb_im.AckMessage{Mid: mid, Seq: seq},
	}
}

// AckNotify 服务端下发给发送者的消息送达通知
type AckNotify struct {
	*pb_im.AckNotify
}

func NewAckNotify(mid int64) AckNotify {
	return AckNotify{
		&pb_im.AckNotify{Mid: mid},
	}
}

type GroupNotify struct {
	*pb_im.GroupNotify
}

func NewGroupNotify(mid, gid int64, seq int64, typ int64, timestamp int64, data interface{}) *GroupNotify {
	notify := &GroupNotify{
		&pb_im.GroupNotify{
			Mid:       mid,
			Gid:       gid,
			Type:      int32(typ),
			Seq:       seq,
			Timestamp: timestamp,
			Data:      &anypb.Any{},
		},
	}
	message, ok := data.(proto.Message)
	if !ok {
		return notify
	}
	_ = notify.Data.MarshalFrom(message)
	return notify
}

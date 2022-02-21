package message

import (
	"errors"
	"fmt"
	"go_im/im/message/pb"
	"go_im/im/message/pb/pb_msg"
	"google.golang.org/protobuf/proto"
)

type Message struct {
	*pb_msg.CommMessage
}

func NewMessage(seq int64, action Action, data interface{}) *Message {
	message := Message{pb.NewMessage(seq, string(action), data)}
	return &message
}

func NewEmptyMessage() *Message {
	return &Message{&pb_msg.CommMessage{Data: &pb_msg.Any{}}}
}

func (m *Message) DeserializeData(v interface{}) error {
	pbMsg, ok := v.(proto.Message)
	if !ok {
		return errors.New("not proto.Message")
	}
	return m.Data.UnmarshalTo(pbMsg)
}

func (m *Message) String() string {
	return fmt.Sprintf("Message{Seq=%d, Action=%s, Data=%v}", m.Seq, m.Action, m.Data)
}

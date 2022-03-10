package message

import (
	"fmt"
	"go_im/im/message/pb"
	"go_im/protobuff/gen/pb_im"
)

type Message struct {
	*pb_im.CommMessage
}

func NewMessage(seq int64, action Action, data interface{}) *Message {
	message := Message{CommMessage: pb.NewMessage(seq, string(action), data)}
	return &message
}

func NewEmptyMessage() *Message {
	return &Message{&pb_im.CommMessage{Data: &pb_im.Any{}}}
}

func (m *Message) DeserializeData(v interface{}) error {
	return DefaultCodec.Decode(m.Data.Value, v)
}

func (m *Message) String() string {
	return fmt.Sprintf("Message{Seq=%d, Action=%s, Data=%v}", m.Seq, m.Action, m.Data)
}

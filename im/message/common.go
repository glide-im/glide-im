package message

import (
	"fmt"
	json2 "go_im/im/message/json"
)

type Message struct {
	json2.CommMessage
}

func NewMessage(seq int64, action Action, data interface{}) *Message {
	return &Message{
		json2.NewMessage(seq, string(action), data),
	}
}

func (m *Message) DeserializeData(v interface{}) error {
	return m.Data.Deserialize(v)
}

func (m *Message) String() string {
	return fmt.Sprintf("Message{Seq=%d, Action=%s, Data=%s}", m.Seq, m.Action, m.Data)
}

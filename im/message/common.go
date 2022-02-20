package message

import (
	"fmt"
	json2 "go_im/im/message/json"
)

type Message struct {
	json2.CommMessage
}

func NewMessage(seq int64, action Action, data interface{}) *Message {
	ret := new(Message)
	ret.Seq = seq
	ret.Action = string(action)
	ret.Data = json2.NewData(data)
	return ret
}

func (m *Message) DeserializeData(v interface{}) error {
	return m.Data.Deserialize(v)
}

func (m *Message) String() string {
	return fmt.Sprintf("Message{Seq=%d, Action=%s, Data=%s}", m.Seq, m.Action)
}

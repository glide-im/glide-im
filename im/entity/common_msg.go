package entity

import (
	"encoding/json"
	"fmt"
)

type Action int32

func (a Action) IsApi() bool {
	return a&MaskActionApi != 0
}

func (a Action) IsMessage() bool {
	return a&MaskActionMessage != 0
}

func (a Action) IsHeartbeat() bool {
	return a == ActionHeartbeat
}

func (a Action) String() string {
	return actionNameMap[a]
}

type Message struct {
	Seq    int64
	Action Action
	Data   string
}

func (m *Message) Deserialize(data []byte) error {
	return json.Unmarshal(data, m)
}

func (m *Message) Serialize() ([]byte, error) {
	return json.Marshal(m)
}

func (m *Message) SetData(v interface{}) error {
	if s, ok := v.(string); ok {
		m.Data = s
		return nil
	}
	if s, ok := v.(error); ok {
		m.Data = s.Error()
		return nil
	}

	b, err := json.Marshal(v)
	if err != nil {
		return err
	}
	m.Data = string(b)
	return nil
}

func (m *Message) DeserializeData(v interface{}) error {
	return json.Unmarshal([]byte(m.Data), v)
}

func (m *Message) String() string {
	n := actionNameMap[m.Action]
	return fmt.Sprintf("Message{Seq=%d, Action=%s, Data=%s}", m.Seq, n, m.Data)
}

func NewMessage(seq int64, action Action, data interface{}) *Message {
	ret := new(Message)
	ret.Seq = seq
	ret.Action = action
	_ = ret.SetData(data)
	return ret
}

package message

import (
	"encoding/json"
	"fmt"
	"strings"
)

const (
	ActionMessage           Action = "message"
	ActionGroupMessage             = "message.group"
	ActionChatMessage              = "message.chat"
	ActionChatMessageRetry         = "message.chat.retry"
	ActionChatMessageResend        = "message.chat.resend"
	ActionCSMessage                = "message.cs"
	ActionMessageFailed            = "message.failed.send"

	ActionAckRequest  = "ack.request"
	ActionAckGroupMsg = "ack.group.msg"
	ActionAckMessage  = "ack.message"
	ActionAckNotify   = "ack.notify"

	ActionApi       = "api"
	ActionHeartbeat = "heartbeat"
	ActionNotify    = "notify"
	ActionFailed    = "failed"
)

type Action string

func (a *Action) Contains(action Action) bool {
	return strings.HasPrefix(string(*a), string(action))
}

type Message struct {
	Ver    int64
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
	return fmt.Sprintf("Message{Seq=%d, Action=%s, Data=%s}", m.Seq, m.Action, m.Data)
}

func NewMessage(seq int64, action Action, data interface{}) *Message {
	ret := new(Message)
	ret.Seq = seq
	ret.Action = action
	_ = ret.SetData(data)
	return ret
}

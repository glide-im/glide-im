package im

import "encoding/json"

type Action int64

const (
	_ Action = iota
	ActionApi
	ActionUserLogin
	ActionUserRegister
	ActionUserGetInfo
	ActionUserEditInfo
	ActionUserLogout

	ActionMessage
	ActionGroupMessage
	ActionChatMessage

	ActionHeartbeat
)

func (a Action) IsApi() bool {
	return a > ActionApi && a < ActionMessage
}

func (a Action) IsMessage() bool {
	return a > ActionMessage && a < ActionHeartbeat
}

func (a Action) IsHeartbeat() bool {
	return a == ActionHeartbeat
}

func (a Action) ActionName() string {
	return ""
}

type Message struct {
	Req    string
	Action Action
	Data   interface{}
}

func DeserializeMessage(data []byte) (*Message, error) {
	m := &Message{}
	err := json.Unmarshal(data, m)
	return m, err
}

func (m *Message) Serialize() ([]byte, error) {
	return json.Marshal(m)
}

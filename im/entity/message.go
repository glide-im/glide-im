package entity

import "encoding/json"

type Action uint64

const (

	// ======================================== api
	MaskActionApi Action = 1 << 40

	ActionUserLogin    = MaskActionApi | 1<<0
	ActionUserRegister = MaskActionApi | 1<<1
	ActionUserGetInfo  = MaskActionApi | 1<<2
	ActionUserEditInfo = MaskActionApi | 1<<3
	ActionUserLogout   = MaskActionApi | 1<<4

	// ======================================== api response
	MaskActionApiResp      = MaskActionApi | MaskActionApi<<1
	ActionUserUnauthorized = MaskActionApiResp | 1<<1

	// ======================================== message
	MaskActionMessage  = 1 << 50
	ActionGroupMessage = MaskActionMessage | 1<<0
	ActionChatMessage  = MaskActionMessage | 1<<1

	// ======================================== heartbeat
	ActionHeartbeat = 1<<60 | 1
)

type Message struct {
	Seq    int64
	Action Action
	Data   []byte
}

func DeserializeMessage(data []byte) (*Message, error) {
	m := &Message{}
	err := json.Unmarshal(data, m)
	return m, err
}

func (m *Message) Serialize() ([]byte, error) {
	return json.Marshal(m)
}

func (m *Message) DeserializeData(v interface{}) error {
	return json.Unmarshal(m.Data, v)
}

func NewErrMessage(seq int64, err error) *Message {
	resp := new(Message)
	resp.Seq = seq
	resp.Data = []byte(err.Error())
	return resp
}

func NewMessage(seq int64, action Action, msg string) *Message {
	ret := new(Message)
	ret.Seq = seq
	ret.Action = action
	ret.Data = []byte(msg)
	return ret
}

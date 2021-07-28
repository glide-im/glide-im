package entity

import (
	"encoding/json"
	"fmt"
)

type Action int32

const (
	ActionAck = 1

	MaskActionApi Action = 1 << 20

	ActionUserLogin       = MaskActionApi | 1
	ActionUserRegister    = MaskActionApi | 2
	ActionUserGetInfo     = MaskActionApi | 3
	ActionUserEditInfo    = MaskActionApi | 4
	ActionUserLogout      = MaskActionApi | 5
	ActionUserChatList    = MaskActionApi | 6
	ActionUserInfo        = MaskActionApi | 7
	ActionUserAuth        = MaskActionApi | 8
	ActionUserRelation    = MaskActionApi | 10
	ActionUserNewChat     = MaskActionApi | 11
	ActionUserChatHistory = MaskActionApi | 12

	ActionOnlineUser = MaskActionApi | 20

	MaskActionMessage  = 1 << 25
	ActionGroupMessage = MaskActionMessage | 1
	ActionChatMessage  = MaskActionMessage | 2

	ActionHeartbeat Action = 1<<30 | 1
)

const (
	MaskRespActionApi          = 1 << 20
	RespActionFailed           = MaskRespActionApi | 1
	RespActionSuccess          = MaskRespActionApi | 2
	RespActionUserUnauthorized = MaskRespActionApi | 3

	MaskRespActionNotify    = 1 << 30
	RespActionGroupRemoved  = MaskRespActionNotify | 1
	RespActionGroupApproval = MaskRespActionNotify | 3
	RespActionGroupApproved = MaskRespActionNotify | 4
	RespActionGroupRefused  = MaskRespActionNotify | 5
	RespActionEcho          = MaskRespActionNotify | 100

	RespActionFriendApproval = MaskRespActionNotify | 6
	RespActionFriendApproved = MaskRespActionNotify | 7
	RespActionFriendRefused  = MaskRespActionNotify | 8
)

var actionNameMap = map[Action]string{
	ActionUserLogin:    "ActionUserLogin",
	ActionUserRegister: "ActionUserRegister",
	ActionUserGetInfo:  "ActionUserGetInfo",
	ActionUserEditInfo: "ActionUserEditInfo",
	ActionUserLogout:   "ActionUserLogout",
	ActionUserChatList: "ActionUserChatList",
	ActionUserInfo:     "ActionUserInfo",
	ActionUserAuth:     "ActionUserAuth",

	ActionUserRelation: "ActionUserRelation",

	MaskActionMessage:  "MaskActionMessage",
	ActionGroupMessage: "ActionGroupMessage",
	ActionChatMessage:  "ActionChatMessage",

	ActionHeartbeat: "ActionHeartbeat",
}

var actionRequestMap map[Action]func() interface{}

type Message struct {
	Seq    int64
	Action Action
	Data   string
}

func DeserializeMessage(data []byte) (*Message, error) {
	m := new(Message)
	err := json.Unmarshal(data, m)
	return m, err
}

func (m *Message) Serialize() ([]byte, error) {
	return json.Marshal(m)
}

func (m *Message) SetData(v interface{}) error {
	if s, ok := v.(string); ok {
		m.Data = s
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

func NewErrMessage(seq int64, err error) *Message {
	resp := new(Message)
	resp.Seq = seq
	resp.Data = err.Error()
	return resp
}

func NewAckMessage(seq int64) *Message {
	resp := new(Message)
	resp.Seq = seq
	resp.Action = ActionAck
	return resp
}

func NewSimpleMessage(seq int64, action Action, msg string) *Message {
	ret := new(Message)
	ret.Seq = seq
	ret.Action = action
	ret.Data = msg
	return ret
}

func NewMessage(seq int64, action Action) *Message {
	ret := new(Message)
	ret.Seq = seq
	ret.Action = action
	return ret
}

func init() {
	actionRequestMap = map[Action]func() interface{}{
		ActionUserLogin:    func() interface{} { return &LoginRequest{} },
		ActionUserRegister: func() interface{} { return &RegisterRequest{} },
		ActionUserGetInfo:  func() interface{} { return &UserInfoRequest{} },

		ActionUserEditInfo:    func() interface{} { return &RegisterRequest{} },
		ActionUserLogout:      nil,
		ActionUserRelation:    nil,
		ActionUserChatList:    nil,
		ActionUserChatHistory: func() interface{} { return &ChatHistoryRequest{} },
		ActionUserInfo:        func() interface{} { return &UserInfoRequest{} },
		ActionUserNewChat:     func() interface{} { return &UserNewChatRequest{} },
	}
}

func NewRequestFromAction(action Action) interface{} {
	fun, ok := actionRequestMap[action]
	if ok && fun != nil {
		return fun()
	}
	return nil
}

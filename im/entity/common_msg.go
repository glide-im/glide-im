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
	ActionUserChatInfo    = MaskActionApi | 13
	ActionUserAddFriend   = MaskActionApi | 14

	ActionOnlineUser = MaskActionApi | 20

	MaskActionGroupApi = MaskActionApi | 1<<21

	ActionGroupCreate       = MaskActionGroupApi | 1
	ActionGroupGetMember    = MaskActionGroupApi | 2
	ActionGroupJoin         = MaskActionGroupApi | 3
	ActionGroupExit         = MaskActionGroupApi | 4
	ActionGroupRemoveMember = MaskActionGroupApi | 5
	ActionGroupInfo         = MaskActionGroupApi | 6
	ActionGroupUpdate       = MaskActionGroupApi | 7
	ActionGroupAddMember    = MaskActionGroupApi | 8

	MaskActionMessage  = 1 << 25
	ActionGroupMessage = MaskActionMessage | 1
	ActionChatMessage  = MaskActionMessage | 2

	MasActionOther         = 1 << 30
	ActionFailed           = MasActionOther | 1
	ActionSuccess          = MasActionOther | 2
	ActionUserUnauthorized = MasActionOther | 3
	ActionNotify           = MasActionOther | 4
	ActionHeartbeat        = MasActionOther | 6
	ActionEcho             = MasActionOther | 100
)

var actionNameMap = map[Action]string{
	ActionUserLogin:       "ActionUserLogin",
	ActionUserRegister:    "ActionUserRegister",
	ActionUserGetInfo:     "ActionUserGetInfo",
	ActionUserEditInfo:    "ActionUserEditInfo",
	ActionUserLogout:      "ActionUserLogout",
	ActionUserChatList:    "ActionUserChatList",
	ActionUserInfo:        "ActionUserInfo",
	ActionUserAuth:        "ActionUserAuth",
	ActionUserRelation:    "ActionUserRelation",
	ActionUserNewChat:     "ActionUserNewChat",
	ActionUserChatHistory: "ActionUserChatHistory",
	ActionUserChatInfo:    "ActionUserChatInfo",
	ActionUserAddFriend:   "",

	ActionGroupRemoveMember: "ActionGroupRemoveMember",
	ActionGroupAddMember:    "ActionGroupAddMember",
	ActionGroupJoin:         "ActionGroupJoin",
	ActionGroupGetMember:    "ActionGroupGetMember",
	ActionGroupExit:         "ActionGroupExit",
	ActionGroupCreate:       "ActionGroupCreate",
	ActionGroupInfo:         "",
	ActionGroupUpdate:       "ActionGroupUpdate",

	ActionOnlineUser: "ActionOnlineUser",

	MaskActionMessage:  "MaskActionMessage",
	ActionGroupMessage: "ActionGroupMessage",
	ActionChatMessage:  "ActionChatMessage",

	ActionHeartbeat:        "ActionHeartbeat",
	ActionFailed:           "ActionFailed",
	ActionSuccess:          "ActionSuccess",
	ActionUserUnauthorized: "ActionUserUnauthorized",
	ActionNotify:           "ActionNotify",
	ActionEcho:             "ActionEcho",
}

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

func init() {
	actionRequestMap = map[Action]func() interface{}{
		ActionUserLogin:       func() interface{} { return &LoginRequest{} },
		ActionUserRegister:    func() interface{} { return &RegisterRequest{} },
		ActionUserGetInfo:     func() interface{} { return &UserInfoRequest{} },
		ActionUserEditInfo:    func() interface{} { return &RegisterRequest{} },
		ActionUserChatHistory: func() interface{} { return &ChatHistoryRequest{} },
		ActionUserChatInfo:    func() interface{} { return &ChatInfoRequest{} },
		ActionUserInfo:        func() interface{} { return &UserInfoRequest{} },
		ActionUserNewChat:     func() interface{} { return &UserNewChatRequest{} },
		ActionUserAddFriend:   func() interface{} { return &AddContacts{} },

		ActionGroupCreate:    func() interface{} { return &CreateGroupRequest{} },
		ActionGroupInfo:      func() interface{} { return &GroupInfoRequest{} },
		ActionGroupJoin:      func() interface{} { return &JoinGroupRequest{} },
		ActionGroupAddMember: func() interface{} { return &AddMemberRequest{} },
		ActionGroupGetMember: func() interface{} { return &GetGroupMemberRequest{} },
		ActionGroupExit:      func() interface{} { return &ExitGroupRequest{} },
	}
}

func NewRequestFromAction(action Action) interface{} {
	fun, ok := actionRequestMap[action]
	if ok && fun != nil {
		return fun()
	}
	return nil
}

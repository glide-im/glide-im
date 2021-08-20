package entity

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

var actionRequestMap map[Action]func() interface{}

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

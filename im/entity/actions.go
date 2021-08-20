package entity

const (
	ActionAck = 1

	ActionUser            = ActionApi + ".user"
	ActionUserAuth        = ActionUser + ".auth"
	ActionUserLogin       = ActionUser + ".login"
	ActionUserLogout      = ActionUser + ".logout"
	ActionUserRegister    = ActionUser + ".register"
	ActionUserGetInfo     = ActionUser + ".info.get"
	ActionUserEditInfo    = ActionUser + ".info.edit"
	ActionUserInfo        = ActionUser + ".info.user"
	ActionUserChatList    = ActionUser + ".chat.list"
	ActionUserNewChat     = ActionUser + ".chat.add"
	ActionUserChatHistory = ActionUser + ".chat.history"
	ActionUserChatInfo    = ActionUser + ".chat.info"
	ActionUserContacts    = ActionUser + ".contacts.get"
	ActionUserAddFriend   = ActionUser + ".contacts.add"
	ActionOnlineUser      = ActionUser + ".online"

	ActionGroup             = ActionApi + ".group"
	ActionGroupCreate       = ActionGroup + ".create"
	ActionGroupJoin         = ActionGroup + ".join"
	ActionGroupExit         = ActionGroup + ".exit"
	ActionGroupUpdate       = ActionGroup + ".update"
	ActionGroupGetMember    = ActionGroup + ".member.get"
	ActionGroupRemoveMember = ActionGroup + ".member.remove"
	ActionGroupAddMember    = ActionGroup + ".member.add"
	ActionGroupInfo         = ActionGroup + ".info.get"

	ActionGroupMessage = ActionMessage + ".group"
	ActionChatMessage  = ActionMessage + ".chat"

	ActionOther            = ""
	ActionFailed           = ""
	ActionSuccess          = ""
	ActionUserUnauthorized = ""
	ActionNotify           = ""
)

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

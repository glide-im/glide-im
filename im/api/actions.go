package api

const (
	ActionAck = 1

	ActionApi = "api"

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

	ActionOther            = ""
	ActionFailed           = ""
	ActionSuccess          = ""
	ActionUserUnauthorized = ""
	ActionNotify           = ""
)

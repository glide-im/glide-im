package http_srv

import (
	"go_im/im/api/auth"
	"go_im/im/api/groups"
	"go_im/im/api/msg"
	"go_im/im/api/user"
)

func initRoute() {

	// TODO 2021-11-15 完成其他 api 的 http 服务

	authApi := auth.AuthApi{}
	post("/api/auth/register", authApi.Register)
	post("/api/auth/logout", authApi.Logout)
	post("/api/auth/signin", authApi.SignIn)

	groupApi := groups.GroupApi{}
	post("/api/group/info", groupApi.GetGroupInfo)
	post("/api/group/members", groupApi.GetGroupMember)
	post("/api/group/create", groupApi.CreateGroup)
	post("/api/group/join", groupApi.JoinGroup)
	post("/api/group/remove", groupApi.RemoveMember)

	userApi := user.UserApi{}
	post("/api/contacts/add", userApi.AddContact)
	post("/api/contacts/list", userApi.GetContactList)
	post("/api/contacts/approval", userApi.ContactApproval)
	post("/api/contacts/del", userApi.DeleteContact)
	post("/api/contacts/update", userApi.UpdateContactRemark)

	post("/api/user/info", userApi.GetUserInfo)
	post("/api/user/profile", userApi.UserProfile)
	post("/api/user/profile/update", userApi.UpdateUserProfile)

	msgApi := msg.MsgApi{}
	post("/api/msg/history", msgApi.GetChatMessageHistory)
	post("/api/msg/recent", msgApi.GetRecentChatMessages)
	post("/api/msg/offline", msgApi.GetOfflineMessage)
	post("/api/msg/offline/ack", msgApi.AckOfflineMessage)

	post("/api/session/recent", msgApi.GetRecentSessions)
	post("/api/session/get", msgApi.GetOrCreateSession)
	post("/api/session/update", msgApi.UpdateSession)
}

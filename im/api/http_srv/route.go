package http_srv

import (
	"go_im/im/api/auth"
	"go_im/im/api/msg"
)

func initRoute() {
	authApi := auth.AuthApi{}

	// TODO 2021-11-15 完成其他 api 的 http 服务
	post("/api/auth/register", authApi.Register)
	post("/api/auth/logout", authApi.Logout)
	post("/api/auth/signin", authApi.SignIn)

	msgApi := msg.MsgApi{}

	post("/api/msg/history", msgApi.GetChatMessageHistory)
	post("/api/msg/recent", msgApi.GetRecentChatMessages)
	post("/api/msg/offline", msgApi.GetOfflineMessage)
	post("/api/msg/offline/ack", msgApi.AckOfflineMessage)

	post("/api/session/recent", msgApi.GetRecentSessions)
	post("/api/session/get", msgApi.GetOrCreateSession)
	post("/api/session/update", msgApi.UpdateSession)
}

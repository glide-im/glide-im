package http_srv

import "go_im/im/api/auth"

func initRoute() {
	authApi := auth.AuthApi{}

	// TODO 2021-11-15 完成其他 api 的 http 服务
	post("/api/auth/register", authApi.Register)
	post("/api/auth/logout", authApi.Logout)
	post("/api/auth/signin", authApi.SignIn)
}

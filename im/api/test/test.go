package test

import (
	route "go_im/im/api/router"
	"go_im/im/client"
	"go_im/im/message"
)

type TestApi struct{}

func (t *TestApi) TestLogin(info *route.Context, request *TestLoginRequest) error {
	return client.SignIn(info.Uid, request.Uid, request.Device)
}

func (t *TestApi) TestSignOut(info *route.Context) error {
	return client.Logout(info.Uid, 2)
}

func (t *TestApi) TestSendMessage(info *route.Context) error {
	return client.EnqueueMessage(2, message.NewMessage(0, "notify.test", "hello world"))
}

package test

import (
	route "github.com/glide-im/glideim/im/api/router"
	"github.com/glide-im/glideim/im/client"
	"github.com/glide-im/glideim/im/message"
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

package test

import (
	route "go_im/im/api/router"
	"go_im/im/client"
	"go_im/im/message"
)

var ResponseHandleFunc func(uid int64, device int64, message *message.Message)

func respond(uid int64, seq int64, action message.Action, data interface{}) {
	resp := message.NewMessage(seq, action, data)
	respondMessage(uid, resp)
}

func respondMessage(uid int64, msg *message.Message) {
	ResponseHandleFunc(uid, 0, msg)
}

type TestApi struct{}

type TestLoginRequest struct {
	Uid    int64
	Device int64
}

func (t *TestApi) TestLogin(info *route.RequestInfo, request *TestLoginRequest) error {
	client.Manager.ClientSignIn(info.Uid, request.Uid, request.Device)
	return nil
}

func (t *TestApi) TestSignOut(info *route.RequestInfo) error {
	client.Manager.ClientLogout(info.Uid, 2)
	return nil
}

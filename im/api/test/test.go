package test

import (
	route "go_im/im/api/router"
	"go_im/im/client"
)

type TestApi struct{}

type TestLoginRequest struct {
	Uid    int64
	Device int64
}

func (t *TestApi) TestLogin(info *route.Context, request *TestLoginRequest) error {
	client.Manager.ClientSignIn(info.Uid, request.Uid, request.Device)
	return nil
}

func (t *TestApi) TestSignOut(info *route.Context) error {
	client.Manager.ClientLogout(info.Uid, 2)
	return nil
}

package test

import (
	route "go_im/im/api/router"
	"go_im/im/client"
)

type TestApi struct{}

func (t *TestApi) TestLogin(info *route.Context, request *TestLoginRequest) error {
	client.SignIn(info.Uid, request.Uid, request.Device)
	return nil
}

func (t *TestApi) TestSignOut(info *route.Context) error {
	client.Logout(info.Uid, 2)
	return nil
}

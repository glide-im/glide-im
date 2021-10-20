package api

import "go_im/im/client"

type TestApi struct{}

type TestLoginRequest struct {
	Uid    int64
	Device int64
}

func (t *TestApi) TestLogin(info *RequestInfo, request *TestLoginRequest) error {
	client.Manager.ClientSignIn(info.Uid, request.Uid, request.Device)
	return nil
}

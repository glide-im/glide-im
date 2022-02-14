package app

import (
	"fmt"
	"go_im/im/api/comm"
	"go_im/im/api/router"
	"go_im/im/dao/appdao"
	"go_im/im/message"
)

type Interface interface {
	Echo(req *route.Context) error
}

type AppApi struct {
}

func (*AppApi) Echo(req *route.Context) error {
	req.Response(message.NewMessage(req.Seq, "api.app.echo", fmt.Sprintf("seq=%d, uid=%d", req.Seq, req.Uid)))
	return nil
}

func (*AppApi) GetReleaseInfo(ctx *route.Context) error {

	info, err := appdao.Impl.GetReleaseInfo()
	if err != nil {
		return err
	}
	ctx.Response(message.NewMessage(0, comm.ActionSuccess, info))
	return nil
}

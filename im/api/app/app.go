package app

import (
	"fmt"
	"github.com/glide-im/glideim/im/api/apidep"
	"github.com/glide-im/glideim/im/api/comm"
	"github.com/glide-im/glideim/im/api/router"
	"github.com/glide-im/glideim/im/client"
	"github.com/glide-im/glideim/im/dao/appdao"
	"github.com/glide-im/glideim/im/message"
	"time"
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

var cacheServerInfo *client.ServerInfo = nil
var cacheInfoExpired = time.Now()

func (a *AppApi) GetServerInfo(ctx *route.Context) error {

	if cacheInfoExpired.After(time.Now()) {
		ctx.ReturnSuccess(cacheServerInfo)
		return nil
	}
	cacheInfoExpired = time.Now().Add(time.Second * 3)

	info := apidep.ClientInterface.GetServerInfo()

	if info == nil {
		ctx.ReturnSuccess(struct{}{})
		return nil
	}

	cacheServerInfo = info
	ctx.ReturnSuccess(info)
	return nil
}

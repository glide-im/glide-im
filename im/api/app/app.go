package app

import (
	"fmt"
	"go_im/im/api/router"
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

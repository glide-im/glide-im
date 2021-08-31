package api

import (
	"fmt"
	"go_im/im/message"
)

type AppApi struct {
}

func (a *AppApi) Echo(req *RequestInfo) error {
	respondMessage(req.Uid, message.NewMessage(req.Seq, "api.app.echo", fmt.Sprintf("seq=%d, uid=%d", req.Seq, req.Uid)))
	return nil
}

package app

import (
	"fmt"
	"go_im/im/api/router"
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

type Interface interface {
	Echo(req *route.RequestInfo) error
}

type AppApi struct {
}

func (AppApi) Echo(req *route.RequestInfo) error {
	respondMessage(req.Uid, message.NewMessage(req.Seq, "api.app.echo", fmt.Sprintf("seq=%d, uid=%d", req.Seq, req.Uid)))
	return nil
}

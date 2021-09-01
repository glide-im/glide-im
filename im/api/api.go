package api

import (
	"go_im/im/client"
	"go_im/im/message"
)

var Impl IApiHandler

type IApiHandler interface {
	Handle(uid int64, message *message.Message)
}

func SetImpl(api IApiHandler) {
	Impl = api
}

func Handle(uid int64, message *message.Message) {
	Impl.Handle(uid, message)
}

func respond(uid int64, seq int64, action message.Action, data interface{}) {
	resp := message.NewMessage(seq, action, data)
	respondMessage(uid, resp)
}

func respondMessage(uid int64, msg *message.Message) {
	client.EnqueueMessage(uid, msg)
}

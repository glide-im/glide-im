package api

import (
	"go_im/im/client"
	"go_im/im/message"
)

var impl Api

type Api interface {
	Handle(uid int64, message *message.Message)
}

func SetImpl(api Api) {
	impl = api
}

func Handle(uid int64, message *message.Message) {
	impl.Handle(uid, message)
}

func respond(uid int64, seq int64, action message.Action, data interface{}) {
	resp := message.NewMessage(seq, action, data)
	respondMessage(uid, resp)
}

func respondMessage(uid int64, msg *message.Message) {
	client.EnqueueMessage(uid, msg)
}

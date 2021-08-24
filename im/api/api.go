package api

import "go_im/im/message"

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

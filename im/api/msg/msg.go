package msg

import (
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
	SyncChatMsgBySeq(msg *route.RequestInfo, request *SyncChatMsgReq) error
}

type MsgApi struct {
}

func (MsgApi) SyncChatMsgBySeq(msg *route.RequestInfo, request *SyncChatMsgReq) error {
	panic("implement me")
}

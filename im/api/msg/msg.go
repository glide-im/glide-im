package msg

import (
	"go_im/im/api/comm"
	route "go_im/im/api/router"
	"go_im/im/dao/msgdao"
	"go_im/im/message"
)

type MsgApi struct {
	*GroupMsgApi
	*ChatMsgApi
}

func (MsgApi) GetMessageID(ctx *route.Context) error {
	id, err := msgdao.GetMessageID()
	if err != nil {
		return comm.NewDbErr(err)
	}
	ctx.Response(message.NewMessage(ctx.Seq, comm.ActionSuccess, MessageIDResponse{id}))
	return nil
}

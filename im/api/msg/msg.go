package msg

import (
	"github.com/glide-im/glideim/im/api/comm"
	route "github.com/glide-im/glideim/im/api/router"
	"github.com/glide-im/glideim/im/dao/msgdao"
	"github.com/glide-im/glideim/im/message"
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

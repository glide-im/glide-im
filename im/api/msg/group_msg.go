package msg

import (
	"go_im/im/api/comm"
	route "go_im/im/api/router"
	"go_im/im/dao/msgdao"
	"go_im/im/message"
)

type GroupMsgApi struct {
}

func (*GroupMsgApi) GetGroupMessageHistory(ctx *route.Context, request *GetGroupMsgRequest) error {

	ms, err := msgdao.GroupMsgDaoImpl.GetGroupMessage(request.Gid, request.Page, 20)
	if err != nil {
		return comm.NewDbErr(err)
	}
	//goland:noinspection GoPreferNilSlice
	resp := []*GroupMessageResponse{}
	for _, m := range ms {
		resp = append(resp, &GroupMessageResponse{
			MID:     m.MID,
			Sender:  m.From,
			Gid:     m.To,
			Type:    m.Type,
			SendAt:  m.SendAt,
			Content: m.Content,
		})
	}
	ctx.Response(message.NewMessage(ctx.Seq, comm.ActionSuccess, resp))
	return nil
}

func (*GroupMsgApi) GetGroupMessageState(ctx *route.Context, request *GetGroupMsgStateRequest) error {

	state, err := msgdao.GroupMsgDaoImpl.GetGroupMessageState(request.Gid)
	if err != nil {
		return comm.NewDbErr(err)
	}
	ctx.Response(message.NewMessage(ctx.Seq, comm.ActionSuccess, GroupMessageStateResponse{state}))
	return nil
}

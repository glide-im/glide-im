package msg

import (
	"go_im/im/api/comm"
	"go_im/im/api/router"
	"go_im/im/dao/msgdao"
	"go_im/im/message"
	"time"
)

type MsgApi struct{}

func (a *MsgApi) UpdateSession(ctx *route.Context, request *SessionRequest) error {
	err := msgdao.SessionDaoImpl.UpdateOrInitSession(ctx.Uid, request.To, time.Now().Unix())
	if err != nil {
		return comm.NewDbErr(err)
	}
	ctx.Response(message.NewMessage(ctx.Seq, comm.ActionSuccess, ""))
	return nil
}

func (*MsgApi) CreateSession(ctx *route.Context, request *SessionRequest) error {
	err := msgdao.SessionDaoImpl.CreateSession(ctx.Uid, request.To)
	if err != nil {
		return comm.NewDbErr(err)
	}
	ctx.Response(message.NewMessage(ctx.Seq, comm.ActionSuccess, ""))
	return nil
}

func (a *MsgApi) GetRecentSessions(ctx *route.Context) error {
	session, err := msgdao.SessionDaoImpl.GetRecentSession(time.Now().Unix() - (time.Hour.Milliseconds()*24*7)/1000)
	if err != nil {
		return comm.NewDbErr(err)
	}
	//goland:noinspection GoPreferNilSlice
	sr := []*SessionResponse{}
	for _, s := range session {
		sr = append(sr, &SessionResponse{
			To:       s.To,
			LastMid:  s.LastMID,
			UpdateAt: s.UpdateAt,
			ReadAt:   s.ReadAt,
		})
	}
	ctx.Response(message.NewMessage(ctx.Seq, comm.ActionSuccess, sr))
	return nil
}

func (*MsgApi) Get(msg *route.Context, request *SyncChatMsgReq) error {
	panic("implement me")
}

func (*MsgApi) SyncChatMsgBySeq2(msg *route.Context, request *SyncChatMsgReq) error {
	panic("implement me")
}

func (*MsgApi) SyncChatMsgBySeq3(msg *route.Context, request *SyncChatMsgReq) error {
	panic("implement me")
}

func (*MsgApi) SyncChatMsgBySeq4(msg *route.Context, request *SyncChatMsgReq) error {
	panic("implement me")
}

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

func (*MsgApi) GetOrCreateSession(ctx *route.Context, request *SessionRequest) error {
	session, err := msgdao.SessionDaoImpl.GetSession(ctx.Uid, request.To)
	if err != nil {
		return comm.NewDbErr(err)
	}
	if session == nil {
		se, err := msgdao.SessionDaoImpl.CreateSession(ctx.Uid, request.To, time.Now().Unix())
		if err != nil {
			return comm.NewDbErr(err)
		}
		session = se
	}

	ctx.Response(message.NewMessage(ctx.Seq, comm.ActionSuccess, session))
	return nil
}

func (a *MsgApi) GetRecentSessions(ctx *route.Context) error {
	week := time.Now().Unix() - (time.Hour.Milliseconds()*24*7)/1000
	session, err := msgdao.SessionDaoImpl.GetRecentSession(ctx.Uid, week)
	if err != nil {
		return comm.NewDbErr(err)
	}
	//goland:noinspection GoPreferNilSlice
	sr := []*SessionResponse{}
	for _, s := range session {
		sr = append(sr, &SessionResponse{
			Uid2:     s.Uid,
			Uid1:     s.Uid2,
			LastMid:  s.LastMID,
			UpdateAt: s.UpdateAt,
		})
	}
	ctx.Response(message.NewMessage(ctx.Seq, comm.ActionSuccess, sr))
	return nil
}

func (*MsgApi) GetRecentChatMessages(msg *route.Context, request *GetRecentMessageRequest) error {

	return nil
}

func (*MsgApi) SyncChatMsgBySeq2(msg *route.Context, request *GetRecentMessageRequest) error {
	panic("implement me")
}

func (*MsgApi) SyncChatMsgBySeq3(msg *route.Context, request *GetRecentMessageRequest) error {
	panic("implement me")
}

func (*MsgApi) SyncChatMsgBySeq4(msg *route.Context, request *GetRecentMessageRequest) error {
	panic("implement me")
}

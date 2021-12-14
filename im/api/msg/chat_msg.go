package msg

import (
	"go_im/im/api/comm"
	"go_im/im/api/router"
	"go_im/im/dao/msgdao"
	"go_im/im/message"
	"go_im/pkg/logger"
	"time"
)

type ChatMsgApi struct{}

//goland:noinspection GoPreferNilSlice
func (*ChatMsgApi) GetChatMessageHistory(ctx *route.Context, request *GetChatHistoryRequest) error {

	ms, err := msgdao.ChatMsgDaoImpl.GetChatMessagesBySession(ctx.Uid, request.Uid, request.Page, 20)
	if err != nil {
		return comm.NewDbErr(err)
	}
	msr := []*MessageResponse{}
	for _, m := range ms {
		msr = append(msr, messageModel2MessageResponse(m))
	}
	ctx.Response(message.NewMessage(ctx.Seq, comm.ActionSuccess, msr))
	return nil
}

//goland:noinspection GoPreferNilSlice
func (*ChatMsgApi) GetRecentMessage(ctx *route.Context) error {
	messages, err := msgdao.ChatMsgDaoImpl.GetRecentChatMessages(ctx.Uid, time.Now().Unix()-int64(time.Hour*3*24))
	if err != nil {
		return comm.NewDbErr(err)
	}
	msr := []*MessageResponse{}
	for _, m := range messages {
		msr = append(msr, messageModel2MessageResponse(m))
	}
	ctx.Response(message.NewMessage(ctx.Seq, comm.ActionSuccess, msr))
	return nil
}

//goland:noinspection GoPreferNilSlice
func (*ChatMsgApi) GetRecentMessageByUser(ctx *route.Context, request *GetRecentMessageRequest) error {
	resp := []RecentMessagesResponse{}
	var e = 0
	for _, i := range request.Uid {
		ms, err := msgdao.ChatMsgDaoImpl.GetChatMessagesBySession(ctx.Uid, i, 0, 20)
		if err != nil {
			logger.E("GetRecentMessageByUser DB error %v", err)
			e++
			continue
		}
		msr := []*MessageResponse{}
		for _, m := range ms {
			msr = append(msr, messageModel2MessageResponse(m))
		}
		resp = append(resp, RecentMessagesResponse{
			Uid:      i,
			Messages: msr,
		})
	}
	if e == len(request.Uid) {
		return errRecentMsgLoadFailed
	}
	ctx.Response(message.NewMessage(ctx.Seq, comm.ActionSuccess, resp))
	return nil
}

func (*ChatMsgApi) AckOfflineMessage(ctx *route.Context, request *AckOfflineMessageRequest) error {
	err := msgdao.DelOfflineMessage(ctx.Uid, request.Mid)
	if err != nil {
		return comm.NewDbErr(err)
	}
	return nil
}

//goland:noinspection GoPreferNilSlice
func (*ChatMsgApi) GetOfflineMessage(ctx *route.Context) error {
	oms, err := msgdao.GetOfflineMessage(ctx.Uid)
	if err != nil {
		return comm.NewDbErr(err)
	}
	var mid = []int64{}
	for _, m := range oms {
		mid = append(mid, m.MID)
	}
	qms, err := msgdao.GetChatMessage(mid...)
	var ms = []*MessageResponse{}
	for _, m := range qms {
		ms = append(ms, messageModel2MessageResponse(m))
	}
	ctx.Response(message.NewMessage(ctx.Seq, comm.ActionSuccess, ms))
	return nil
}

func messageModel2MessageResponse(m *msgdao.ChatMessage) *MessageResponse {
	return &MessageResponse{
		MID:      m.MID,
		CliSeq:   m.CliSeq,
		From:     m.From,
		To:       m.To,
		Type:     m.Type,
		SendAt:   m.SendAt,
		CreateAt: m.CreateAt,
		Content:  m.Content,
	}
}

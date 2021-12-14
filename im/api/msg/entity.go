package msg

import "go_im/im/dao/msgdao"

type MessageResponse struct {
	MID      int64
	CliSeq   int64
	From     int64
	To       int64
	Type     int
	SendAt   int64
	CreateAt int64
	Content  string
}

type GroupMessageResponse struct {
	MID     int64
	Sender  int64
	Gid     int64
	Type    int
	SendAt  int64
	Content string
}

type GroupMessageStateResponse struct {
	*msgdao.GroupMessageState
}

type SessionRequest struct {
	To int64
}

type SessionResponse struct {
	Uid1     int64
	Uid2     int64
	LastMid  int64
	UpdateAt int64
}

type GetRecentMessageRequest struct {
	Uid []int64
}

type RecentMessagesResponse struct {
	Uid      int64
	Messages []*MessageResponse
}

type GetChatHistoryRequest struct {
	Uid  int64
	Page int
}

type AckOfflineMessageRequest struct {
	Mid []int64
}

type ChatHistoryRequest struct {
	Cid  int64
	Time int64
	Type int8
}

type GetGroupMsgRequest struct {
	Gid  int64
	Page int
}

type GetGroupMsgStateRequest struct {
	Gid int64
}

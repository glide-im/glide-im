package msg

import "go_im/im/dao/msgdao"

type MessageResponse struct {
	Mid      int64
	CliSeq   int64
	From     int64
	To       int64
	Type     int
	SendAt   int64
	CreateAt int64
	Content  string
}

type GroupMessageResponse struct {
	Mid     int64
	Sender  int64
	Gid     int64
	Seq     int64
	Type    int
	SendAt  int64
	Content string
}

type GroupMessageStateResponse struct {
	*msgdao.GroupMessageState
}

type ReadMessageRequest struct {
	To int64
}

type SessionRequest struct {
	To int64
}

type SessionResponse struct {
	Uid1     int64
	Uid2     int64
	Unread   int64
	LastMid  int64
	UpdateAt int64
	CreateAt int64
}

type RecentChatMessageRequest struct {
	Uid int64
}

type RecentMessageRequest struct {
	Uid []int64
}

type RecentMessagesResponse struct {
	Uid      int64
	Messages []*MessageResponse
}

type ChatHistoryRequest struct {
	Uid       int64
	BeforeMid int64
}

type AckOfflineMessageRequest struct {
	Mid []int64
}

type GroupMessageRequest struct {
	Mid []int64
}

type RecentGroupMessageRequest struct {
	Gid int64
}

type GroupMsgHistoryRequest struct {
	Gid       int64
	BeforeSeq int64
}

type GroupMsgStateRequest struct {
	Gid int64
}

type MessageIDResponse struct {
	Mid int64
}

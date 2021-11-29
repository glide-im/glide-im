package msg

type SessionRequest struct {
	To int64
}

type SessionResponse struct {
	To       int64
	LastMid  int64
	UpdateAt int64
	ReadAt   int64
}

type SyncChatMsgReq struct {
}

type ChatHistoryRequest struct {
	Cid  int64
	Time int64
	Type int8
}

type ChatInfoRequest struct {
	Cid int64
}

package msg

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

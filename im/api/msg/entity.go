package msg

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

type ChatHistoryRequest struct {
	Cid  int64
	Time int64
	Type int8
}

type ChatInfoRequest struct {
	Cid int64
}

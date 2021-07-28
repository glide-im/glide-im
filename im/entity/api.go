package entity

type LoginRequest struct {
	Device   int64  `json:"device"`
	Account  string `json:"account"`
	Password string `json:"password"`
}

type AuthRequest struct {
	Token    string `json:"token"`
	DeviceId int64  `json:"device_id"`
	Uid      int64  `json:"uid"`
}

type RegisterRequest struct {
	Account  string `json:"account"`
	Password string `json:"password"`
}

// AuthorResponse login or register result
type AuthorResponse struct {
	Token string
	Uid   int64
}

type UserInfoRequest struct {
	Uid []int64
}

type UserInfoResponse struct {
	Uid      int64
	Nickname string
	Avatar   string
}

type UserInfoListResponse struct {
	UserInfo []*UserInfoResponse
}

type UserNewChatRequest struct {
	Id   uint64
	Type int8
}

type RelationResponse struct {
	Groups  []uint64
	Friends []int64
}

type ChatHistoryRequest struct {
	Cid  uint64
	Time int64
	Type int8
}

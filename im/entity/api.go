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
	Token string `json:"token"`
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

type RelationResponse struct {
	Groups  []uint64
	Friends []int64
}

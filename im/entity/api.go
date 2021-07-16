package entity

type LoginRequest struct {
	Device   int64  `json:"device"`
	Username string `json:"username"`
	Password string `json:"password"`
}

type AuthRequest struct {
	Token    string `json:"token"`
	DeviceId int64  `json:"device_id"`
}

type RegisterRequest struct {
	Username string `json:"username"`
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

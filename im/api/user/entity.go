package user

type InfoRequest struct {
	Uid []int64
}

type InfoResponse struct {
	Uid      int64
	Nickname string
	Account  string
	Avatar   string
}

type InfoListResponse struct {
	UserInfo []*InfoResponse
}

type UpdateProfileRequest struct {
}
type ContactResponse struct {
	Id     int64
	Type   int8
	Remark string
}

type AddContacts struct {
	Uid    int64
	Remark string
}

type DeleteContactsRequest struct {
	Uid []int64
}

type UpdateRemarkRequest struct {
	Uid    int64
	Remark string
}

type ContactApproval struct {
	Uid     int64
	Agree   bool
	Comment string
}

type OnlineUser struct {
	Uid    int64
	Before int64
}

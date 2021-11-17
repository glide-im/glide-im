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

type NewChatRequest struct {
	Id   int64
	Type int8
}

type ContactResponse struct {
	Friends []*InfoResponse
	Groups  []interface{}
}
type AddContacts struct {
	Uid    int64
	Remark string
}

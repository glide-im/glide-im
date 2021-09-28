package api

import "go_im/im/dao"

type LoginRequest struct {
	Device   int64
	Account  string
	Password string
}

type AuthRequest struct {
	Token    string
	DeviceId int64
}

type RegisterRequest struct {
	Account  string
	Password string
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
	Account  string
	Avatar   string
}

type UserInfoListResponse struct {
	UserInfo []*UserInfoResponse
}

type UserNewChatRequest struct {
	Id   int64
	Type int8
}

type ContactResponse struct {
	Friends []*UserInfoResponse
	Groups  []*GroupResponse
}

type ChatHistoryRequest struct {
	Cid  int64
	Time int64
	Type int8
}

type ChatInfoRequest struct {
	Cid int64
}

type GroupInfoRequest struct {
	Gid []int64
}

type CreateGroupRequest struct {
	Name   string
	Member []int64
}

type GroupResponse struct {
	dao.Group
	Members []*dao.GroupMember
}

type GroupAddMemberResponse struct {
	Gid     int64
	Members []*dao.GroupMember
}

type AddedGroupResponse struct {
	Group *dao.Group
	UcId  int64
}

type JoinGroupRequest struct {
	Gid int64
}

type ExitGroupRequest struct {
	Gid int64
}

type GetGroupMemberRequest struct {
	Gid int64
}

type GroupMemberResponse struct {
	Uid        int64
	Nickname   string
	RemarkName string
	Type       int32
	Online     bool
	Mute       bool
}

type AddMemberRequest struct {
	Gid int64
	Uid []int64
}

type AddContacts struct {
	Uid    int64
	Remark string
}

type RemoveMemberRequest struct {
	Gid int64
	Uid []int64
}

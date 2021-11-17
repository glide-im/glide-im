package groups

import (
	"go_im/im/dao/groupdao"
)

type GroupInfoRequest struct {
	Gid []int64
}

type CreateGroupRequest struct {
	Name   string
	Member []int64
}

type GroupResponse struct {
	groupdao.Group
	Members []*groupdao.GroupMember
}

type GroupAddMemberResponse struct {
	Gid     int64
	Members []*groupdao.GroupMember
}

type AddedGroupResponse struct {
	Group *groupdao.Group
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

type RemoveMemberRequest struct {
	Gid int64
	Uid []int64
}

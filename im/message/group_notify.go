package message

import (
	"go_im/protobuff/gen/pb_im"
)

const (
	GroupNotifyTypeMemberAdded       = 1
	GroupNotifyTypeMemberRemoved     = 2
	GroupNotifyTypeMemberSetAdmin    = 3
	GroupNotifyTypeMemberRemoveAdmin = 4
	GroupNotifyTypeMemberSetMute     = 5
	GroupNotifyTypeMemberRemoveMute  = 6
)

type GroupNotifyMemberAdded struct {
	pb_im.GroupNotifyMemberAdded
}

func NewGroupNotifyAdded(uid []int64) GroupNotifyMemberAdded {
	return GroupNotifyMemberAdded{
		pb_im.GroupNotifyMemberAdded{
			Uid: uid,
		},
	}
}

type GroupNotifyMemberRemove struct {
	pb_im.GroupNotifyMemberRemove
}

func NewGroupNotifyRemove(uid []int64) GroupNotifyMemberRemove {
	return GroupNotifyMemberRemove{
		pb_im.GroupNotifyMemberRemove{
			Uid: uid,
		},
	}
}

package message

import (
	"go_im/im/message/pb/pb_msg"
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
	pb_msg.GroupNotifyMemberAdded
}

func NewGroupNotifyAdded(uid []int64) GroupNotifyMemberAdded {
	return GroupNotifyMemberAdded{
		pb_msg.GroupNotifyMemberAdded{
			Uid: uid,
		},
	}
}

type GroupNotifyMemberRemove struct {
	pb_msg.GroupNotifyMemberRemove
}

func NewGroupNotifyRemove(uid []int64) GroupNotifyMemberRemove {
	return GroupNotifyMemberRemove{
		pb_msg.GroupNotifyMemberRemove{
			Uid: uid,
		},
	}
}

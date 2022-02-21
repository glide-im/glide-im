package json

const (
	GroupNotifyTypeMemberAdded       = 1
	GroupNotifyTypeMemberRemoved     = 2
	GroupNotifyTypeMemberSetAdmin    = 3
	GroupNotifyTypeMemberRemoveAdmin = 4
	GroupNotifyTypeMemberSetMute     = 5
	GroupNotifyTypeMemberRemoveMute  = 6
)

type GroupNotifyMemberAdded struct {
	Uid []int64
}

type GroupNotifyMemberRemove struct {
	Uid []int64
}

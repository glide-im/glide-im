package group

import (
	"go_im/im/dao"
	"go_im/im/message"
)

var Manager IGroupManager

type IGroupManager interface {
	PutMember(gid int64, mb map[int64]int32)

	RemoveMember(gid int64, uid ...int64) error

	UserOnline(uid, gid int64)

	UserOffline(uid, gid int64)

	AddGroup(group *dao.Group, owner int64)

	DispatchNotifyMessage(uid int64, gid int64, message *message.Message)

	DispatchMessage(uid int64, message *message.Message) error
}

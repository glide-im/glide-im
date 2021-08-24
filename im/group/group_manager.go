package group

import (
	"go_im/im/dao"
	"go_im/im/message"
)

var Manager IGroupManager

type IGroupManager interface {
	PutMember(gid int64, mb *dao.GroupMember)

	UnsubscribeGroup(uid, gid int64)

	RemoveMember(gid int64, uid ...int64) error

	GetMembers(gid int64) ([]*dao.GroupMember, error)

	AddGroup(group *Group)

	GetGroup(gid int64) *Group

	DispatchNotifyMessage(uid int64, gid int64, message *message.Message)

	DispatchMessage(uid int64, message *message.Message) error
}

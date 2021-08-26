package group

import (
	"go_im/im/dao"
	"go_im/im/message"
)

var Manager IGroupManager

type IGroupManager interface {
	PutMember(gid int64, mb *dao.GroupMember)

	//UnsubscribeGroup(uid, gid int64)

	RemoveMember(gid int64, uid ...int64) error

	GetMembers(gid int64) ([]*dao.GroupMember, error)

	AddGroup(group *dao.Group, cid int64, owner *dao.GroupMember)

	GetGroup(gid int64) *dao.Group

	GetGroupCid(gid int64) int64

	HasMember(gid int64, uid int64) bool

	DispatchNotifyMessage(uid int64, gid int64, message *message.Message)

	DispatchMessage(uid int64, message *message.Message) error
}

package apidep

import (
	"go_im/im/client"
	"go_im/im/group"
	"go_im/im/message"
)

// ClientManager 客户端连接相关接口
var ClientManager ClientManagerInterface = client.Manager

// GroupManager 群相关接口
var GroupManager GroupManagerInterface = group.Manager

type ClientManagerInterface interface {
	ClientSignIn(oldUid int64, uid int64, device int64)
	ClientLogout(uid int64, device int64)
	EnqueueMessage(uid int64, device int64, message *message.Message)
	IsDeviceOnline(uid, device int64) bool
	IsOnline(uid int64) bool
	AllClient() []int64
}

type GroupManagerInterface interface {
	PutMember(gid int64, mb map[int64]int32) error
	RemoveMember(gid int64, uid ...int64) error
	AddGroup(gid int64) error
	RemoveGroup(gid int64) error
	ChangeStatus(gid int64, status int64) error
	DispatchNotifyMessage(gid int64, message *message.Message) error
}

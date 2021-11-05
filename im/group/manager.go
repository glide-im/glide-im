package group

import (
	"go_im/im/message"
)

// Manager 群相关操作入口
var Manager IGroupManager

type IGroupManager interface {
	// PutMember 给指定添加群成员, mb 为 uid-type, 用户ID-成员类型
	PutMember(gid int64, mb map[int64]int32)

	// RemoveMember 移除群成员
	RemoveMember(gid int64, uid ...int64) error

	// AddGroup 添加群, 创建
	AddGroup(gid int64)

	// RemoveGroup 移除群, 解散
	RemoveGroup(gid int64)

	// ChangeStatus 设置群状态, 禁言等
	ChangeStatus(gid int64, status int64)

	// DispatchNotifyMessage 发送通知消息
	DispatchNotifyMessage(gid int64, message *message.Message)

	// DispatchMessage 发送聊天消息
	DispatchMessage(gid int64, message *message.GroupMessage)
}

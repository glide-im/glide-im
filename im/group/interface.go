package group

import (
	"go_im/im/client"
	"go_im/im/message"
)

const (
	_ = iota
	FlagMemberAdd
	FlagMemberDel
	FlagMemberOnline
	FlagMemberOffline
	FlagMemberMuted
	FlagMemberSetAdmin
	FlagMemberCancelAdmin
)

const (
	_ = iota
	FlagGroupCreate
	FlagGroupDissolve
	FlagGroupMute
	FlagGroupCancelMute
)

type MessageHandler func(uid int64, device int64, message *message.Message) error

type MemberUpdate struct {
	Uid  int64
	Flag int64

	Extra interface{}
}

type Update struct {
	Flag int64

	Extra interface{}
}

type Interface interface {
	// UpdateMember 更新群成员
	UpdateMember(gid int64, update []MemberUpdate) error

	// UpdateGroup 更新群
	UpdateGroup(gid int64, update Update) error

	// DispatchNotifyMessage 发送通知消息
	DispatchNotifyMessage(gid int64, message *message.GroupNotify) error

	// DispatchMessage 发送聊天消息
	DispatchMessage(gid int64, action message.Action, message *message.ChatMessage) error
}

// manager 群相关操作入口
var manager Interface = NewDefaultManager()

var enqueueMessage MessageHandler = client.EnqueueMessageToDevice

func SetMessageHandler(handler MessageHandler) {
	enqueueMessage = handler
}

func SetInterfaceImpl(i Interface) {
	manager = i
}

func UpdateMember(gid int64, update []MemberUpdate) error {
	return manager.UpdateMember(gid, update)
}

// UpdateGroup 更新群
func UpdateGroup(gid int64, update Update) error {
	return manager.UpdateGroup(gid, update)
}

// DispatchNotifyMessage 发送通知消息
func DispatchNotifyMessage(gid int64, message *message.GroupNotify) error {
	return manager.DispatchNotifyMessage(gid, message)
}

// DispatchMessage 发送聊天消息
func DispatchMessage(gid int64, msg *message.ChatMessage) error {
	return manager.DispatchMessage(gid, message.ActionChatMessage, msg)
}

func DispatchRecallMessage(gid int64, msg *message.ChatMessage) error {
	return manager.DispatchMessage(gid, message.ActionGroupMessageRecall, msg)
}

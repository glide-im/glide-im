package group

import (
	"github.com/glide-im/glideim/im/client"
	"github.com/glide-im/glideim/im/message"
)

const (
	FlagMemberOnline      int64 = 1 << 62
	FlagMemberOffline           = 1 << 61
	FlagMemberMuted             = 1 << 1
	FlagMemberTypeAdmin         = 1 << 2
	FlagMemberTypeGeneral       = 1 << 3
)

const (
	FlagGroupCreate     int64 = 1
	FlagGroupDissolve         = 2
	FlagGroupMute             = 3
	FlagGroupCancelMute       = 4
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

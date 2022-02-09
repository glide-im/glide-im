package group

import "go_im/im/message"

func UpdateMember(gid int64, update []MemberUpdate) error {
	return Manager.UpdateMember(gid, update)
}

// UpdateGroup 更新群
func UpdateGroup(gid int64, update Update) error {
	return Manager.UpdateGroup(gid, update)
}

// DispatchNotifyMessage 发送通知消息
func DispatchNotifyMessage(gid int64, message *message.GroupNotify) error {
	return Manager.DispatchNotifyMessage(gid, message)
}

// DispatchMessage 发送聊天消息
func DispatchMessage(gid int64, message *message.UpChatMessage) error {
	return Manager.DispatchMessage(gid, message)
}

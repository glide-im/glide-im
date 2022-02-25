package client

import "go_im/im/message"

// MessageHandleFunc 所有客户端消息都传递到该函数处理
var MessageHandleFunc func(from int64, device int64, message *message.Message) = nil

// Manager 客户端管理入口
var Manager IClientManager = NewDefaultManager()

func SignIn(oldUid int64, uid int64, device int64) {
	Manager.ClientSignIn(oldUid, uid, device)
}
func Logout(uid int64, device int64) {
	Manager.ClientLogout(uid, device)
}
func IsDeviceOnline(uid, device int64) bool {
	return false
}
func IsOnline(uid int64) bool {
	return false
}
func AllClient() []int64 {
	return []int64{}
}

// EnqueueMessage Manager.EnqueueMessage 的快捷方法, 预留一个位置对消息入队列进行一些预处理
func EnqueueMessage(uid int64, message *message.Message) {
	//
	Manager.EnqueueMessage(uid, 0, message)
}

func EnqueueMessageToDevice(uid int64, device int64, message *message.Message) {
	Manager.EnqueueMessage(uid, device, message)
}

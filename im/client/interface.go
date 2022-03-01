package client

import "go_im/im/message"

type Interface interface {
	ClientSignIn(oldUid int64, uid int64, device int64)

	ClientLogout(uid int64, device int64)

	EnqueueMessage(uid int64, device int64, message *message.Message)
}

// MessageHandleFunc 所有客户端消息都传递到该函数处理
var MessageHandleFunc func(from int64, device int64, message *message.Message) = nil

// Manager 客户端管理入口
var manager Interface = NewDefaultManager()

func SignIn(oldUid int64, uid int64, device int64) error {
	manager.ClientSignIn(oldUid, uid, device)
	return nil
}
func Logout(uid int64, device int64) error {
	manager.ClientLogout(uid, device)
	return nil
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
func EnqueueMessage(uid int64, message *message.Message) error {
	//
	manager.EnqueueMessage(uid, 0, message)
	return nil
}

func EnqueueMessageToDevice(uid int64, device int64, message *message.Message) error {
	manager.EnqueueMessage(uid, device, message)
	return nil
}

func SetInterfaceImpl(i Interface) {
	manager = i
}

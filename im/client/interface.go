package client

import "go_im/im/message"

type Interface interface {
	ClientSignIn(oldUid int64, uid int64, device int64) error

	ClientLogout(uid int64, device int64) error

	EnqueueMessage(uid int64, device int64, message *message.Message) error
}

type MessageHandler func(from int64, device int64, message *message.Message) error

// messageHandleFunc 所有客户端消息都传递到该函数处理
var messageHandleFunc MessageHandler = nil

// Manager 客户端管理入口
var manager Interface = NewDefaultManager()

func SignIn(oldUid int64, uid int64, device int64) error {
	return manager.ClientSignIn(oldUid, uid, device)
}
func Logout(uid int64, device int64) error {
	return manager.ClientLogout(uid, device)
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
	return manager.EnqueueMessage(uid, 0, message)
}

func EnqueueMessageToDevice(uid int64, device int64, message *message.Message) error {
	return manager.EnqueueMessage(uid, device, message)
}

func SetInterfaceImpl(i Interface) {
	manager = i
}

func SetMessageHandler(handler MessageHandler) {
	messageHandleFunc = handler
}

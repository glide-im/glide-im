package client

import (
	"go_im/im/conn"
	"go_im/im/message"
)

var Manager IClientManager

var MessageHandleFunc MessageHandler = nil

type MessageHandler func(from int64, message *message.Message) error

type IClientManager interface {
	ClientConnected(conn conn.Connection) int64

	AddClient(uid int64, cs IClient)

	ClientSignIn(oldUid int64, uid int64, device int64)

	ClientLogout(uid int64)

	EnqueueMessage(uid int64, message *message.Message)

	AllClient() []int64
}

func EnqueueMessage(uid int64, message *message.Message) {
	Manager.EnqueueMessage(uid, message)
}

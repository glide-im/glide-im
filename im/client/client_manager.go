package client

import (
	"go_im/im/conn"
	"go_im/im/message"
)

var Manager IClientManager

type IClientManager interface {
	ClientConnected(conn conn.Connection) int64

	ClientSignIn(oldUid int64, uid int64, device int64)

	UserLogout(uid int64)

	HandleMessage(from int64, message *message.Message) error

	EnqueueMessage(uid int64, message *message.Message)

	AllClient() []int64
}

func EnqueueMessage(uid int64, message *message.Message) {
	Manager.EnqueueMessage(uid, message)
}

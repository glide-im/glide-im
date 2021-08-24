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

	DispatchMessage(from int64, message *message.Message) error

	Api(from int64, message *message.Message)

	EnqueueMessage(uid int64, message *message.Message)

	IsOnline(uid int64) bool

	AllClient() []int64

	Update()
}

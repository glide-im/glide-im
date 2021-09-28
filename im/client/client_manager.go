package client

import (
	"go_im/im/conn"
	"go_im/im/message"
)

// Manager 客户端管理入口
var Manager IClientManager

// IClientManager 管理所有客户端的创建, 消息派发, 退出等
type IClientManager interface {

	// ClientConnected 当一个用户连接建立后, 由该方法创建 IClient 实例 Client 并管理该连接, 返回该由连接创建客户端的标识 id
	// 返回的标识 id 是一个临时 id, 后续连接认证后会改变
	ClientConnected(conn conn.Connection) int64

	// AddClient 用于手段创建一个 IClient, 方便自定义临时 uid 以及其他的 IClient 实现
	AddClient(uid int64, cs IClient)

	// ClientSignIn 给一个已存在的客户端设置一个新的 id, 若 uid 已存在, 则新增一个 device 共享这个 id
	ClientSignIn(oldUid int64, uid int64, device int64)

	// ClientLogout 指定 uid 的客户端退出
	ClientLogout(uid int64)

	// EnqueueMessage 尝试将消息放入指定 uid 的客户端
	EnqueueMessage(uid int64, message *message.Message)

	// AllClient 返回所有的客户端 id
	AllClient() []int64
}

// EnqueueMessage Manager.EnqueueMessage 的快捷方法, 预留一个位置对消息入队列进行一些预处理
func EnqueueMessage(uid int64, message *message.Message) {
	Manager.EnqueueMessage(uid, message)
}

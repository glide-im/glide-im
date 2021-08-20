package im

import (
	"go_im/im/conn"
	"go_im/im/dao"
)

func Run() {

	dao.Init()
	wsServer := conn.NewWsServer(nil)
	wsServer.Handler(func(conn conn.Connection) {
		ClientManager.ClientConnected(conn)
	})
	wsServer.Run()
}

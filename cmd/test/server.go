package test

import (
	"go_im/im/client"
	"go_im/im/conn"
	"go_im/im/dao"
	"go_im/im/group"
	"go_im/pkg/db"
	"time"
)

func Run() {
	db.Init()
	dao.Init()

	var server conn.Server

	op := &conn.WsServerOptions{
		ReadTimeout:  10 * time.Minute,
		WriteTimeout: 10 * time.Minute,
	}
	server = conn.NewWsServer(op)

	server.SetConnHandler(func(conn conn.Connection) {
		client.Manager.ClientConnected(conn)
	})

	client.Manager = client.NewDefaultManager()
	manager := group.NewDefaultManager()
	group.Manager = manager
	manager.Init()

	err := server.Run("0.0.0.0", 8080)
	if err != nil {
		panic(err)
	}
}

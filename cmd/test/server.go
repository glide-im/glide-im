package test

import (
	"go_im/im/api"
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
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}
	server = conn.NewWsServer(op)

	server.SetConnHandler(func(conn conn.Connection) {
		client.Manager.ClientConnected(conn)
	})

	api.MessageHandleFunc = client.EnqueueMessage
	api.SetHandler(api.NewApiRouter())
	client.Manager = client.NewClientManager()
	manager := group.NewGroupManager()
	group.Manager = manager
	manager.Init()

	err := server.Run("0.0.0.0", 8080)
	if err != nil {
		panic(err)
	}
}

package main

import (
	"go_im/im/client"
	"go_im/im/conn"
	"go_im/im/dao"
	"go_im/im/dao/msgdao"
	"go_im/im/group"
	"go_im/pkg/db"
	"time"
)

func RunTestServer() {
	db.Init()
	msgdao.MockChatMsg(time.Millisecond * 10)
	msgdao.MockCommDao()
	dao.Init()

	var server conn.Server

	op := &conn.WsServerOptions{
		ReadTimeout:  10 * time.Minute,
		WriteTimeout: 10 * time.Minute,
	}
	server = conn.NewWsServer(op)

	cm := client.NewDefaultManager()
	server.SetConnHandler(func(conn conn.Connection) {
		cm.ClientConnected(conn)
	})

	client.SetInterfaceImpl(cm)
	manager := group.NewDefaultManager()
	group.SetInterfaceImpl(manager)
	manager.Init()

	err := server.Run("0.0.0.0", 8080)
	if err != nil {
		panic(err)
	}
}

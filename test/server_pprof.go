package main

import (
	"github.com/glide-im/glideim/im/client"
	"github.com/glide-im/glideim/im/conn"
	"github.com/glide-im/glideim/im/dao"
	"github.com/glide-im/glideim/im/dao/msgdao"
	"github.com/glide-im/glideim/im/group"
	"github.com/glide-im/glideim/pkg/db"
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

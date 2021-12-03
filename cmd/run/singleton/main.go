package main

import (
	"go_im/im/client"
	"go_im/im/conn"
	"go_im/im/dao"
	"go_im/im/group"
	"go_im/im/messaging"
	"go_im/pkg/db"
	"time"
)

func main() {
	Run()
}

func Run() {
	db.Init()
	dao.Init()

	var server conn.Server

	op := &conn.WsServerOptions{
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}
	messaging.Init()
	server = conn.NewWsServer(op)

	server.SetConnHandler(func(conn conn.Connection) {
		client.Manager.ClientConnected(conn)
	})

	group.Manager.(*group.DefaultManager).Init()

	err := server.Run("0.0.0.0", 8080)
	if err != nil {
		panic(err)
	}
}

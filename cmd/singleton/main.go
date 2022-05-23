package main

import (
	"go_im/config"
	"go_im/im/api"
	"go_im/im/client"
	"go_im/im/conn"
	"go_im/im/dao"
	"go_im/im/group"
	"go_im/im/messaging"
	"go_im/pkg/db"
	"sync"
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
		ReadTimeout:  60 * time.Second,
		WriteTimeout: 60 * time.Second,
	}
	client.SetMessageHandler(messaging.HandleMessage)
	server = conn.NewWsServer(op)

	cm := client.NewDefaultManager()
	server.SetConnHandler(func(conn conn.Connection) {
		cm.ClientConnected(conn)
	})

	client.SetInterfaceImpl(cm)

	manager := group.NewDefaultManager()
	group.SetInterfaceImpl(manager)
	manager.Init()

	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		addr := config.ApiHttp.Addr
		port := config.ApiHttp.Port
		err := api.RunHttpServer(addr, port)
		if err != nil {
			panic(err)
		}
		wg.Done()
	}()

	wg.Add(1)
	go func() {
		addr := config.WsServer.Addr
		port := config.WsServer.Port
		err := server.Run(addr, port)
		if err != nil {
			panic(err)
		}
		wg.Done()
	}()

	wg.Wait()
}

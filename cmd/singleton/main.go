package main

import (
	"github.com/glide-im/glideim/config"
	"github.com/glide-im/glideim/im/api"
	"github.com/glide-im/glideim/im/client"
	"github.com/glide-im/glideim/im/conn"
	"github.com/glide-im/glideim/im/dao"
	"github.com/glide-im/glideim/im/group"
	"github.com/glide-im/glideim/im/messaging"
	"github.com/glide-im/glideim/pkg/db"
	"sync"
	"time"
)

func main() {

	err := config.Load()
	if err != nil {
		panic(err)
	}

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
